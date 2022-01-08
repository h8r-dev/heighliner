package main

import (
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/server"
)

func main() {
	var port int

	logger.Init()

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start running server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := server.New(port)
			return s.Start()
		},
	}

	cmd.Flags().IntVar(&port, "port", 3000, "The number of the port to serve the http APIs.")

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
