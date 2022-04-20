package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a logger with the corresponding format and level
func New() *zap.Logger {
	debugEnabler := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	infoEnabler := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	warnEnabler := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	errorEnabler := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	// High-priority output should go to standard error, and low-priority
	// output should go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	coreTree := zapcore.NewTee(
		zapcore.NewCore(getZapEncoder(), consoleDebugging, debugEnabler),
		zapcore.NewCore(getZapEncoder(), consoleDebugging, infoEnabler),
		zapcore.NewCore(getZapEncoder(), consoleDebugging, warnEnabler),
		zapcore.NewCore(getZapEncoder(), consoleErrors, errorEnabler),
	)

	logger := zap.New(coreTree, zap.AddCaller(), zap.AddStacktrace(errorEnabler))
	return logger
}

// getZapEncodingConfig returns the configuration of zap encoder.
func getZapEncodingConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
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
