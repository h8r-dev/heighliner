package kube

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// NewConfig returns a config
func NewConfig() (*rest.Config, error) {
	restConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	return restConfig, nil
}

// NewClientSet returns a clientset
func NewClientSet() (*kubernetes.Clientset, error) {
	restConfig, err := NewConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

// NewDynamicClient returns a dynamicClient
func NewDynamicClient() (dynamic.Interface, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}
	dcli, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return dcli, nil
}
