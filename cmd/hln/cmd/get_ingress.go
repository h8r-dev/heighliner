package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	watchtools "k8s.io/client-go/tools/watch"

	"github.com/h8r-dev/heighliner/internal/k8sfactory"
)

const (
	defaultNS  = "ingress-nginx"
	defaultSVC = "ingress-nginx-controller"
	igLabel    = "app.kubernetes.io/component=controller"
	defaultIP  = "127.0.0.1" // This IP is for local kind and minikube
)

type getIngressOptions struct {
	genericclioptions.IOStreams
}

func newGetIngressCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &getIngressOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "ingress",
		Short: "Get ingress IP address",
		Args:  cobra.NoArgs,
		RunE:  o.runGetIngress,
	}

	return cmd
}

func (o *getIngressOptions) runGetIngress(cmd *cobra.Command, args []string) error {
	defaultTimeout := 90 * time.Second
	if err := waitForIGController(defaultNS, igLabel, defaultTimeout); err != nil {
		return err
	}
	ip, err := getIngressIP(defaultNS, defaultSVC)
	if err != nil {
		return err
	}
	fmt.Fprintln(o.Out, ip)
	return nil
}

func waitForIGController(namespace, label string, timeout time.Duration) error {
	endTime := time.Now().Add(timeout)
	cs, err := k8sfactory.GetDefaultClientSet()
	if err != nil {
		return err
	}
	for {
		timeout = time.Until(endTime)
		ctx, cancel := watchtools.ContextWithOptionalTimeout(context.Background(), timeout)
		objWatch, err := cs.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{
			LabelSelector: label,
		})
		if err != nil {
			return err
		}
		_, err = watchtools.UntilWithoutRetry(ctx, objWatch, podIsReady)
		cancel()
		switch {
		case err == nil:
			return nil
		case errors.Is(err, watchtools.ErrWatchClosed):
			continue
		case errors.Is(err, wait.ErrWaitTimeout):
			return errors.New("watch time out")
		default:
			return err
		}
	}
}

// IsDeleted returns true if the object is deleted. It prints any errors it encounters.
func podIsReady(event watch.Event) (bool, error) {
	pod, ok := event.Object.(*v1.Pod)
	if !ok {
		return false, errors.New("non-pod resource")
	}
	ready := true
	for _, status := range pod.Status.ContainerStatuses {
		if !status.Ready {
			ready = false
		}
	}
	return ready, nil
}

func getIngressIP(namespace, svcName string) (string, error) {
	cs, err := k8sfactory.GetDefaultClientSet()
	if err != nil {
		return "", err
	}
	ctx := context.TODO()
	igsvc, err := cs.CoreV1().Services(namespace).Get(ctx, svcName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	igs := igsvc.Status.LoadBalancer.Ingress
	var igip string
	if len(igs) > 0 {
		igip = igs[0].IP
	} else {
		igip = defaultIP
	}
	return igip, err
}
