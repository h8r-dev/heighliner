package k8sfactory

import (
	"k8s.io/client-go/kubernetes"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

var (
	defaultFactory cmdutil.Factory
)

// GetDefaultFactory for cluster operations.
func GetDefaultFactory() cmdutil.Factory {
	if defaultFactory == nil {
		return k8sutil.NewFactory(k8sutil.GetKubeConfigPath())
	}
	return defaultFactory
}

// GetDefaultClientSet for cluster operations.
func GetDefaultClientSet() (*kubernetes.Clientset, error) {
	return GetDefaultFactory().KubernetesClientSet()
}
