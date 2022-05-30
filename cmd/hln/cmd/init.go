package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/terraform"
)

func newInitCmd(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize dependent tools and services",
	}
	cmd.RunE = func(c *cobra.Command, args []string) error {
		return checkAndInstall(streams)
	}
	// Shadow the root PersistentPreRun
	cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {}

	return cmd
}

func checkAndInstall(streams genericclioptions.IOStreams) error {
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
	return nil
	// nhctlCli, err := nhctl.NewDefaultClient(streams)
	// if err != nil {
	// 	return err
	// }
	// return nhctlCli.CheckAndInstall()
}
