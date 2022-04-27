package kube

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"

	"github.com/h8r-dev/heighliner/pkg/logger"
)

// TODO move the contents and delete this package.

var argoAppResource = schema.GroupVersionResource{
	Group:    "argoproj.io",
	Version:  "v1alpha1",
	Resource: "applications",
}

const argoCDFinalizerRaw = `{"metadata": {"finalizers": ["resources-finalizer.argocd.argoproj.io"]}}`

func patchFinalizerAndDelete(ctx context.Context,
	client dynamic.Interface,
	namespace, name string,
	streams genericclioptions.IOStreams) error {
	lg := logger.New(streams)
	argoApp := client.Resource(argoAppResource).Namespace(namespace)
	_, err := argoApp.Patch(ctx, name, types.MergePatchType, []byte(argoCDFinalizerRaw), metav1.PatchOptions{})
	if err != nil {
		return err
	}
	lg.Info(fmt.Sprintf("patch finalizer to app %s", name))
	return argoApp.Delete(ctx, name, metav1.DeleteOptions{})
}

// DeleteArgoCDApps cleans up all argocd apps in the specified namespace.
func DeleteArgoCDApps(ctx context.Context,
	client dynamic.Interface,
	namespace string,
	streams genericclioptions.IOStreams) error {
	lg := logger.New(streams)
	argoApp := client.Resource(argoAppResource).Namespace(namespace)
	list, err := argoApp.List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, app := range list.Items {
		appname := app.GetName()
		err := patchFinalizerAndDelete(ctx, client, namespace, appname, streams)
		if err != nil {
			return err
		}
		lg.Info(fmt.Sprintf("app %s deleted", appname))
	}
	return nil
}
