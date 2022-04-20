package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// New creates a logger with the corresponding format and level
func New(streams genericclioptions.IOStreams) *zap.Logger {
	// High-priority output should go to standard error, and low-priority
	// output should go to standard out.
	consoleDebugging := zapcore.AddSync(streams.Out)
	consoleErrors := zapcore.AddSync(streams.ErrOut)
	cores := []zapcore.Core{}
	ll := viper.GetString("log-level")
	// Set default level to "info"
	if ll != "debug" && ll != "warn" && ll != "error" {
		ll = "info"
	}
	switch ll {
	case "debug":
		cores = append(cores,
			zapcore.NewCore(getZapEncoder(),
				consoleDebugging,
				zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zapcore.DebugLevel
				})))
		fallthrough
	case "info":
		cores = append(cores,
			zapcore.NewCore(getZapEncoder(),
				consoleDebugging,
				zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zapcore.InfoLevel
				})))
		fallthrough
	case "warn":
		cores = append(cores,
			zapcore.NewCore(getZapEncoder(),
				consoleDebugging,
				zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev == zapcore.WarnLevel
				})))
		fallthrough
	case "error":
		cores = append(cores,
			zapcore.NewCore(getZapEncoder(),
				consoleErrors,
				zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
					return lev >= zapcore.ErrorLevel
				})))
	}
	return zap.New(zapcore.NewTee(cores...), zap.AddCaller())
}

// getZapEncodingConfig returns the configuration of zap encoder.
func getZapEncodingConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// getZapEncoder returns the zap encoder.
func getZapEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(getZapEncodingConfig())
}
