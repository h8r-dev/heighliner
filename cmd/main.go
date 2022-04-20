package main

import (
	"github.com/h8r-dev/heighliner/pkg/cmd"
)

func main() {
	cmd.Execute(cmd.NewRootCmd())
}
