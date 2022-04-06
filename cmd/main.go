package main

import (
	"github.com/h8r-dev/heighliner/pkg/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	cmd.Execute(rootCmd)
}
