package main

import (
	"github.com/rs/zerolog/log"

	"github.com/h8r-dev/heighliner/pkg/commands/clientcmd"
	"github.com/h8r-dev/heighliner/pkg/logger"
)

func main() {
	logger.Init()

	err := clientcmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to execute client command")
	}
}
