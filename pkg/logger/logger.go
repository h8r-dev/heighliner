package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// New creates a logger with the corresponding format and level
func New() zerolog.Logger {
	logger := zerolog.
		New(os.Stderr).
		With().
		Timestamp().
		Logger()

	if !jsonLogs() {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		logger = logger.With().Timestamp().Caller().Logger()
	}

	level := viper.GetString("log-level")
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	return logger.Level(lvl)
}

func jsonLogs() bool {
	switch f := viper.GetString("log-format"); f {
	case "json":
		return true
	case "plain":
		return false
	case "auto":
		return !term.IsTerminal(int(os.Stdout.Fd()))
	default:
		fmt.Fprintf(os.Stderr, "invalid --log-format %q\n", f)
		os.Exit(1)
	}
	return false
}
