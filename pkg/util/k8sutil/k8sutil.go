package k8sutil

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/util/homedir"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"os"
	"path/filepath"
)

// NewFactory can generate many kinds of k8s client.
func NewFactory(kubeconfigPath string) cmdutil.Factory {
	configFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	configFlags.KubeConfig = &kubeconfigPath
	return cmdutil.NewFactory(configFlags)
}

func GetKubeConfigPath() string {
	path, ok := os.LookupEnv("KUBECONFIG")
	if ok && path != "" {
		return path
	}
	return filepath.Join(homedir.HomeDir(), ".kube", "config")
}
