package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init inits the logger
func Init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
