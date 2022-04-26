package k8sutil

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// NewFactory can generate many kinds of k8s client.
func NewFactory(kubeconfigPath string) cmdutil.Factory {
	configFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	configFlags.KubeConfig = &kubeconfigPath
	return cmdutil.NewFactory(configFlags)
}
