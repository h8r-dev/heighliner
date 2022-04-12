package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir"

	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

// LogsOptions controls the behavior of logs command.
type LogsOptions struct {
	Namespace string
	Pod       string

	// PodLogOptions
	Follow    bool
	Container string

	ClientSet *kubernetes.Clientset
}

func newLogsCmd() *cobra.Command {
	o := &LogsOptions{}

	cmd := &cobra.Command{
		Use:   "logs [POD]",
		Short: "Print the logs for a container in a pod",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			kubeconfigPath := cmd.Flags().Lookup("kubeconfig").Value.String()
			o.ClientSet, err = k8sutil.MakeKubeClient(kubeconfigPath)
			if err != nil {
				return fmt.Errorf("failed to make kube client: %w", err)
			}
			o.Pod = args[0]
			return o.getPodLogs()
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

func (o *LogsOptions) addFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&o.Follow, "follow", "f", o.Follow, "Specify if the logs should be streamed.")
	cmd.Flags().StringVarP(&o.Container, "container", "c", o.Container, "Print the logs of this container")
	cmd.Flags().StringVar(&o.Namespace, "namespace", "default", "Specify the namespace")
}

func (o *LogsOptions) getPodLogs() error {
	request := o.ClientSet.CoreV1().Pods(o.Namespace).GetLogs(o.Pod, &v1.PodLogOptions{
		Container: o.Container,
		Follow:    o.Follow,
	})
	return DefaultConsumeRequest(request, os.Stdout)
}

// DefaultConsumeRequest reads the data from request and writes into
// the out writer. It buffers data from requests until the newline or io.EOF
// occurs in the data, so it doesn't interleave logs sub-line
// when running concurrently.
func DefaultConsumeRequest(request rest.ResponseWrapper, out io.Writer) error {
	readCloser, err := request.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer func() {
		if err := readCloser.Close(); err != nil {
			panic(err)
		}
	}()

	r := bufio.NewReader(readCloser)
	for {
		bytes, err := r.ReadBytes('\n')
		if _, err := out.Write(bytes); err != nil {
			return err
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
}
