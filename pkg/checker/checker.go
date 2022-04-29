package checker

import (
	"github.com/fatih/color"
	"go.uber.org/zap"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/nhctl"
	"github.com/h8r-dev/heighliner/pkg/terraform"
)

// TODO move this two functions and delete this package.

// PreCheck just check the infras and print tips.
func PreCheck(streams genericclioptions.IOStreams) error {
	prompt := "please run hln init"
	lg := logger.New(streams)
	ioDiscard := genericclioptions.NewTestIOStreamsDiscard()
	daggerCli, err := dagger.NewDefaultClient(ioDiscard)
	if err != nil {
		return err
	}
	if err := daggerCli.Check(); err != nil {
		lg.Warn(color.HiYellowString(prompt),
			zap.NamedError("warn", err))
	}
	nhctlCli, err := nhctl.NewDefaultClient(ioDiscard)
	if err != nil {
		return err
	}
	if err := nhctlCli.Check(); err != nil {
		lg.Warn(color.HiYellowString(prompt),
			zap.NamedError("warn", err))
	}
	tfCli, err := terraform.NewDefaultClient(ioDiscard)
	if err != nil {
		return err
	}
	if err := tfCli.Check(); err != nil {
		lg.Warn(color.HiYellowString(prompt),
			zap.NamedError("warn", err))
	}
	return nil
}

// Check will install the infras.
func Check(streams genericclioptions.IOStreams) error {
	daggerCli, err := dagger.NewDefaultClient(streams)
	if err != nil {
		return err
	}
	if err := daggerCli.CheckAndInstall(); err != nil {
		return err
	}
	tfCli, err := terraform.NewDefaultClient(streams)
	if err != nil {
		return err
	}
	if err := tfCli.CheckAndInstall(); err != nil {
		return err
	}
	nhctlCli, err := nhctl.NewDefaultClient(streams)
	if err != nil {
		return err
	}
	return nhctlCli.CheckAndInstall()
}
