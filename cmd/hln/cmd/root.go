package cmd

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/terraform"
)

const greetBanner = `
██╗  ██╗███████╗██╗ ██████╗ ██╗  ██╗██╗     ██╗███╗   ██╗███████╗██████╗ 
██║  ██║██╔════╝██║██╔════╝ ██║  ██║██║     ██║████╗  ██║██╔════╝██╔══██╗
███████║█████╗  ██║██║  ███╗███████║██║     ██║██╔██╗ ██║█████╗  ██████╔╝
██╔══██║██╔══╝  ██║██║   ██║██╔══██║██║     ██║██║╚██╗██║██╔══╝  ██╔══██╗
██║  ██║███████╗██║╚██████╔╝██║  ██║███████╗██║██║ ╚████║███████╗██║  ██║
╚═╝  ╚═╝╚══════╝╚═╝ ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝
`

type configure struct {
	genericclioptions.IOStreams
}

var cfg = configure{
	IOStreams: genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	},
}

// NewRootCmd creates and returns the root command of hln
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hln",
		Short: "Heighliner: Cloud native best practices to build and deploy your applications",
		Long:  greetBanner,
		// Don't print usage message.
		SilenceUsage: true,
		// Logger will print the errors.
		SilenceErrors: true,
	}

	cmd.PersistentPreRunE = func(c *cobra.Command, args []string) error {
		return preCheck(cfg.IOStreams)
	}

	cmd.AddCommand(
		newListCmd(cfg.IOStreams),
		newVersionCmd(cfg.IOStreams),
		newUpCmd(cfg.IOStreams),
		newDownCmd(cfg.IOStreams),
		newStatusCmd(cfg.IOStreams),
		newLogsCmd(cfg.IOStreams),
		newMetricsCmd(cfg.IOStreams),
		newInitCmd(cfg.IOStreams),
		newDomainMappingCmd(cfg.IOStreams),
		newShowCmd(cfg.IOStreams),
	)

	cmd.PersistentFlags().String("log-format", "plain", "Log format (auto, plain, json)")
	cmd.PersistentFlags().StringP("log-level", "l", "info", "Log level")
	// Bind persistent flags to viper
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		log.Fatal().Err(err).Msg("failed to bind flags")
	}

	// Hide 'completion' command
	cmd.CompletionOptions.HiddenDefaultCmd = true

	viper.SetEnvPrefix("hln")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	return cmd
}

func preCheck(streams genericclioptions.IOStreams) error {
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
	// nhctlCli, err := nhctl.NewDefaultClient(ioDiscard)
	// if err != nil {
	// 	return err
	// }
	// if err := nhctlCli.Check(); err != nil {
	// 	lg.Warn(color.HiYellowString(prompt),
	// 		zap.NamedError("warn", err))
	// }
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

// Execute executes the root command with context
func Execute(rootCmd *cobra.Command) {
	var (
		ctx = appcontext.Context()
		lg  = logger.New(cfg.IOStreams)
	)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		lg.Error(err.Error())
		os.Exit(1)
	}
}
