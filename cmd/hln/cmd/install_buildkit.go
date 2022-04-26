package cmd

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"

	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

type InstallBuildKitOptions struct {
	Namespace string

	Kubecli *kubernetes.Clientset
}

func (o *InstallBuildKitOptions) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Namespace, "namespace", "default", "Specify the namespace")
}

func newInstallBuildKitCmd() *cobra.Command {
	o := &InstallBuildKitOptions{}

	cmd := &cobra.Command{
		Use:   "install-buildkit",
		Short: "Install buildkit to cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			kubeconfigPath := cmd.Flags().Lookup("kubeconfig").Value.String()
			o.Kubecli, err = k8sutil.NewFactory(kubeconfigPath).KubernetesClientSet()
			if err != nil {
				return fmt.Errorf("failed to make kube client: %w", err)
			}
			return o.Run()
		},
	}
	if home := homedir.HomeDir(); home != "" {
		cmd.Flags().StringP("kubeconfig", "", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		cmd.Flags().StringP("kubeconfig", "", "", "(optional) absolute path to the kubeconfig file")
	}
	o.addFlags(cmd)
	return cmd
}

func (o *InstallBuildKitOptions) Run() error {
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
			ProbeHandler:        corev1.ProbeHandler{Exec: &corev1.ExecAction{[]string{"buildctl", "debug", "workers"}}},
			InitialDelaySeconds: 5,
			PeriodSeconds:       10,
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler:        corev1.ProbeHandler{Exec: &corev1.ExecAction{[]string{"buildctl", "debug", "workers"}}},
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

	f, _ := fields.ParseSelector(fmt.Sprintf("metadata.name=%s", buildKitName))
	watchlist := cache.NewListWatchFromClient(
		o.Kubecli.AppsV1().RESTClient(),
		"deployments",
		o.Namespace,
		f, //fields.Everything()
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
		return errors.New(fmt.Sprintf("Waiting %s to be ready failed: timeout for 5 minutes", buildKitName))
	}
}
