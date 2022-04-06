package main

import (
	"log"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/h8r-dev/heighliner/pkg/cmd"
)

func main() {
	rootPath := "../heighliner-website/docs/07-cli/hln/commands"
	if len(os.Args) > 1 {
		rootPath = os.Args[1]
	}
	hln := cmd.NewRootCmd()
	err := doc.GenMarkdownTree(hln, rootPath)
	if err != nil {
		log.Fatal(err)
	}
}
