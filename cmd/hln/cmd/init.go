package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/fluxcd/pkg/untar"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/cache"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/hlnpath"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/terraform"
	"github.com/h8r-dev/heighliner/pkg/util/getter"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

const infraSrc = "https://stack.h8r.io/infra.tar.gz"

type initOptions struct {
	WithoutDashboard bool

	genericclioptions.IOStreams
}

func (o *initOptions) BindFlags(f *pflag.FlagSet) {
	f.BoolVar(&o.WithoutDashboard, "without-dashboard", false, "Don't install hln dashboard")
}

func newInitCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &initOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize dependent tools and services",
	}
	cmd.RunE = func(c *cobra.Command, args []string) error {
		if err := checkAndInstall(streams); err != nil {
			return err
		}
		return o.initInfrasForCluster()
	}
	o.BindFlags(cmd.Flags())

	// Shadow the root PersistentPreRun
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {}

	return cmd
}

func (o *initOptions) initInfrasForCluster() error {
	if err := installBuildKit(); err != nil {
		return err
	}
	if err := runForward(o.IOStreams); err != nil {
		return err
	}
	if err := o.runInfraStack(); err != nil {
		return err
	}
	st, err := getStateInSpecificBackend()
	if err != nil {
		return err
	}
	infra, err := st.LoadInfra()
	if err != nil {
		return fmt.Errorf("failed to load infrastructure info: %w", err)
	}
	fmt.Fprintf(o.Out, "\nPlease visit %s to see the dashboard\n", color.CyanString(infra.Dashboard.Ingress))
	fmt.Fprintf(o.Out, "\tUsername: %s\n\tPassword: %s\n", infra.Dashboard.Credentials.Username, infra.Dashboard.Credentials.Password)
	documentationLink := "https://heighliner.dev/docs/getting_started/installation"
	fmt.Fprintf(o.Out, "\nSee the documentation for more information: %s\n", documentationLink)
	return nil
}

func (o *initOptions) runInfraStack() error {
	if o.WithoutDashboard {
		if err := os.Setenv("HLN_WITHOUT_DASHBOARD", "true"); err != nil {
			return err
		}
	}

	kc, ok := os.LookupEnv("KUBECONFIG")
	if !ok {
		kc = k8sutil.GetKubeConfigPath()
	}
	if err := os.Setenv("KUBECONFIG", kc); err != nil {
		return err
	}

	infraPath := hlnpath.CachePath("infrastructure", "infra")
	if err := os.RemoveAll(infraPath); err != nil {
		return err
	}
	src := infraSrc
	dst := filepath.Dir(infraPath)
	tarName := "infra.tar.gz"
	if err := getter.Get(os.Stdout, getter.NewRequest(src, dst, tarName)); err != nil {
		return fmt.Errorf("failed to pull infra source: %w", err)
	}
	tarFile := filepath.Join(dst, tarName)
	data, err := os.ReadFile(tarFile)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)
	if _, err := untar.Untar(buf, dst); err != nil {
		return err
	}
	if err := os.Remove(tarFile); err != nil {
		return err
	}

	cli, err := dagger.NewClient(
		viper.GetString("log-format"),
		viper.GetString("log-level"),
		o.IOStreams,
	)
	if err != nil {
		return err
	}
	return cli.Do(&dagger.ActionOptions{
		Name: "up",
		Dir:  infraPath,
		Plan: "./plan",
	})
}

func checkAndInstall(streams genericclioptions.IOStreams) error {
	errCh := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		daggerCli, err := dagger.NewDefaultClient(streams)
		if err != nil {
			errCh <- err
			return
		}
		if err := daggerCli.CheckAndInstall(); err != nil {
			errCh <- err
			return
		}
	}()
	go func() {
		defer wg.Done()
		tfCli, err := terraform.NewDefaultClient(streams)
		if err != nil {
			errCh <- err
			return
		}
		if err := tfCli.CheckAndInstall(); err != nil {
			errCh <- err
			return
		}
	}()
	wg.Wait()
	return nil
	// nhctlCli, err := nhctl.NewDefaultClient(streams)
	// if err != nil {
	// 	return err
	// }
	// return nhctlCli.CheckAndInstall()
}

func installBuildKit() error {
	client, err := k8sutil.NewFactory(k8sutil.GetKubeConfigPath()).KubernetesClientSet()
	if err != nil {
		return fmt.Errorf("failed to make kube client: %w", err)
	}
	// Create namespace if not exist
	_, err = client.CoreV1().Namespaces().Get(context.TODO(), state.HeighlinerNs, metav1.GetOptions{})
	if err != nil {
		if !k8serr.IsNotFound(err) {
			return err
		}
		var ns corev1.Namespace
		ns.Name = state.HeighlinerNs
		_, err = client.CoreV1().Namespaces().Create(context.TODO(), &ns, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	_, err = client.AppsV1().Deployments(state.HeighlinerNs).Get(context.TODO(), buildKitName, metav1.GetOptions{})
	if err == nil {
		fmt.Println(buildKitName + " has already been installed, skip it")
		return nil
	}

	buildKitLabels := map[string]string{"app": "buildkitd"}
	var buildKitDeploy v1.Deployment
	privileged := true

	buildKitDeploy.Name = buildKitName
	buildKitDeploy.Labels = buildKitLabels
	buildKitDeploy.Spec.Selector = &metav1.LabelSelector{MatchLabels: buildKitLabels}
	buildKitDeploy.Spec.Template.Labels = buildKitLabels
	buildKitDeploy.Spec.Template.Spec.Containers = []corev1.Container{{
		Name:  buildKitName,
		Image: "moby/buildkit:master",
		Args:  []string{"--addr", "unix:///run/buildkit/buildkitd.sock", "--addr", "tcp://0.0.0.0:1234"},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler:        corev1.ProbeHandler{Exec: &corev1.ExecAction{Command: []string{"buildctl", "debug", "workers"}}},
			InitialDelaySeconds: 5,
			PeriodSeconds:       30,
			FailureThreshold:    10,
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler:        corev1.ProbeHandler{Exec: &corev1.ExecAction{Command: []string{"buildctl", "debug", "workers"}}},
			InitialDelaySeconds: 5,
			PeriodSeconds:       30,
			FailureThreshold:    10,
		},
		SecurityContext: &corev1.SecurityContext{Privileged: &privileged},
		Ports:           []corev1.ContainerPort{{ContainerPort: 1234}},
	}}
	_, err = client.AppsV1().Deployments(state.HeighlinerNs).Create(context.TODO(), &buildKitDeploy, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Deployment %s created\n", buildKitName)

	f, _ := fields.ParseSelector(fmt.Sprintf("metadata.name=%s", buildKitName))
	watchlist := cache.NewListWatchFromClient(
		client.AppsV1().RESTClient(),
		"deployments",
		state.HeighlinerNs,
		f,
	)

	stop := make(chan struct{})
	defer close(stop)
	readyCh := make(chan struct{})
	_, controller := cache.NewInformer(
		// also take a look at NewSharedIndexInformer
		watchlist,
		&v1.Deployment{},
		0,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				obj, ok := newObj.(*v1.Deployment)
				if !ok {
					fmt.Printf("expected a *apps.Deployment, got %T\n", obj)
					return
				}

				for _, c := range obj.Status.Conditions {
					if c.Type == v1.DeploymentAvailable && c.Status == "True" {
						readyCh <- struct{}{}
						return
					}
				}
			},
		},
	)
	go controller.Run(stop)

	fmt.Printf("Waiting %s to be ready...\n", buildKitName)
	select {
	case <-readyCh:
		fmt.Printf("%s is ready!\n", buildKitName)
		return nil
	case <-time.Tick(5 * time.Minute):
		return fmt.Errorf("waiting %s to be ready failed: timeout for 5 minutes ", buildKitName)
	}
}
