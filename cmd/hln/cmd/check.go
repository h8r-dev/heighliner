package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/homedir"

	"github.com/h8r-dev/heighliner/pkg/checker"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

func newCheckCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &checkOptions{}
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check if the infrastructures are available",
	}
	cmd.RunE = func(c *cobra.Command, args []string) error {
		err := checker.Check(streams)
		if err != nil {
			return err
		}
		if o.InstallBuildKit {
			kubeconfigPath := cmd.Flags().Lookup("kubeconfig").Value.String()
			o.Kubecli, err = k8sutil.NewFactory(kubeconfigPath).KubernetesClientSet()
			if err != nil {
				return fmt.Errorf("failed to make kube client: %w", err)
			}
			return o.Run()
		}
		return nil
	}
	// Shadow the root PersistentPreRun
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {}

	if home := homedir.HomeDir(); home != "" {
		cmd.Flags().StringP("kubeconfig", "", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		cmd.Flags().StringP("kubeconfig", "", "", "(optional) absolute path to the kubeconfig file")
	}
	o.addFlags(cmd)
	return cmd
}

type checkOptions struct {
	Namespace       string
	InstallBuildKit bool

	Kubecli *kubernetes.Clientset
}

func (o *checkOptions) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Namespace, "namespace", "default", "Specify the namespace")
	cmd.Flags().BoolVar(&o.InstallBuildKit, "install-buildkit", false, "Install buildkit to cluster")
}

func (o *checkOptions) Run() error {
	buildKitName := "buildkitd"
	buildKitLabels := map[string]string{"app": "buildkitd"}
	_, err := o.Kubecli.AppsV1().Deployments(o.Namespace).Get(context.TODO(), buildKitName, metav1.GetOptions{})
	if err == nil {
		fmt.Println(buildKitName + " has already been installed, skip it")
		return nil
	}

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
			PeriodSeconds:       10,
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler:        corev1.ProbeHandler{Exec: &corev1.ExecAction{Command: []string{"buildctl", "debug", "workers"}}},
			InitialDelaySeconds: 5,
			PeriodSeconds:       10,
		},
		SecurityContext: &corev1.SecurityContext{Privileged: &privileged},
		Ports:           []corev1.ContainerPort{{ContainerPort: 1234}},
	}}
	_, err = o.Kubecli.AppsV1().Deployments(o.Namespace).Create(context.TODO(), &buildKitDeploy, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Deployment %s craeted\n", buildKitName)

	f, _ := fields.ParseSelector(fmt.Sprintf("metadata.name=%s", buildKitName))
	watchlist := cache.NewListWatchFromClient(
		o.Kubecli.AppsV1().RESTClient(),
		"deployments",
		o.Namespace,
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
