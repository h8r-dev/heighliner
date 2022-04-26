package checker

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
	"github.com/h8r-dev/heighliner/pkg/util/nhctl"
)

// PreFlight just check the infras and print tips.
func PreFlight(streams genericclioptions.IOStreams) error {
	lg := logger.New(streams)
	dc, err := dagger.NewDefaultClient(
		genericclioptions.NewTestIOStreamsDiscard())
	if err != nil {
		return err
	}
	if err := dc.Check(); err != nil {
		lg.Warn(color.HiYellowString("please run hln check"),
			zap.NamedError("warn", err))
	}
	return nil
}

// Check will install the infras.
func Check(streams genericclioptions.IOStreams) error {
	f := k8sutil.NewFactory("")
	clientSet, err := f.KubernetesClientSet()
	if err != nil {
		return err
	}
	list, err := clientSet.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, t := range list.Items {
		fmt.Println(t.GetName())
	}

	dc, err := dagger.NewDefaultClient(streams)
	if err != nil {
		return err
	}
	if err := dc.CheckAndInstall(); err != nil {
		return err
	}
	if err := nhctl.Check(); err != nil {
		return err
	}
	return util.Exec(streams, nhctl.GetPath(), "version")
}
