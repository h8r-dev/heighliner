package commands

import (
	"strings"

	"github.com/moby/buildkit/util/appcontext"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/h8r-dev/heighliner/pkg/logger"
)

var (
	rootCmd = &cobra.Command{
		Use:   "hln",
		Short: "Heighliner: Cloud native best practices to build and deploy your applications",
	}
)

func init() {
	rootCmd.PersistentFlags().String("log-format", "auto", "Log format (auto, plain, tty, json)")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Log level")

	rootCmd.AddCommand(
		cleanCmd,
		stackCmd,
		newCmd,
		inputCmd,
		upCmd,
	)

	// Hide 'completion' command
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		panic(err)
	}
	viper.SetEnvPrefix("hln")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

// Execute executes the root command.
func Execute() {
	var (
		ctx = appcontext.Context()
		lg  = logger.New()
	)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		lg.Fatal().Err(err).Msg("failed to execute command")
	}
}
