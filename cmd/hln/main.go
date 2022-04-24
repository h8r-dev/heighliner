package main

import (
	"github.com/h8r-dev/heighliner/cmd/hln/cmd"
)

func main() {
	cmd.Execute(cmd.NewRootCmd())
}
