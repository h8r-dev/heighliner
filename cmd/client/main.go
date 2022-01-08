package main

import (
	"github.com/spf13/cobra"

	"github.com/h8r-dev/heighliner/pkg/logger"
)

func main() {
	logger.Init()

	cmd := &cobra.Command{
		Use:   "hln",
		Short: "Heighliner client",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. read stack
			// 2. init + new plan + input
			// 3. dagger up
			return nil
		},
	}

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
