package cmd

import (
	"os"
	"strings"

	"github.com/moby/buildkit/util/appcontext"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/h8r-dev/heighliner/pkg/checker"
	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
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
		return checker.PreCheck(cfg.IOStreams)
	}

	cmd.AddCommand(
		newListCmd(),
		newVersionCmd(),
		newUpCmd(cfg.IOStreams),
		newDownCmd(cfg.IOStreams),
		newStatusCmd(),
		newLogsCmd(),
		newMetricsCmd(),
		newInitCmd(cfg.IOStreams),
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

var (
	fact cmdutil.Factory
)

func getDefaultFactory() cmdutil.Factory {
	if fact == nil {
		return k8sutil.NewFactory(k8sutil.GetKubeConfigPath())
	}
	return fact
}

func getDefaultClientSet() (*kubernetes.Clientset, error) {
	return getDefaultFactory().KubernetesClientSet()
}
