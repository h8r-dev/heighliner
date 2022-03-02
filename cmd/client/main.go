package main

import (
	"github.com/rs/zerolog/log"

	"github.com/h8r-dev/heighliner/pkg/commands"
	"github.com/h8r-dev/heighliner/pkg/logger"
)

func main() {
	logger.Init()

	err := commands.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to execute client command")
	}
}
