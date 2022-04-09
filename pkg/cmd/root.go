package cmd

import (
	"strings"

	"github.com/moby/buildkit/util/appcontext"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/h8r-dev/heighliner/pkg/logger"
)

// NewRootCmd creates and returns the root command of hln
func NewRootCmd() *cobra.Command {
	var err error

	rootCmd := &cobra.Command{
		Use:   "hln",
		Short: "Heighliner: Cloud native best practices to build and deploy your applications",
	}

	rootCmd.AddCommand(
		newListCmd(),
		newVersionCmd(),
		newUpCmd(),
		newTestCmd(),
	)

	rootCmd.PersistentFlags().String("log-format", "auto", "Log format (auto, plain, json)")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Log level")
	// Bind flags to viper
	if err = viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal().Err(err).Msg("failed to bind flags")
	}

	// Hide 'completion' command
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	viper.SetEnvPrefix("hln")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	return rootCmd
}

// Execute executes the root command with context
func Execute(rootCmd *cobra.Command) {
	var (
		ctx = appcontext.Context()
		lg  = logger.New()
	)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		lg.Fatal().Err(err).Msg("failed to execute command")
	}
}
