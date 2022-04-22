package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra/doc"

	"github.com/h8r-dev/heighliner/pkg/cmd"
)

func main() {
	docPath := ""
	if len(os.Args) > 1 {
		docPath = os.Args[1]
	}
	if docPath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		docPath = filepath.Join(pwd, "docs", "commands")
	}
	if err := os.RemoveAll(docPath); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(docPath, 0755); err != nil {
		panic(err)
	}
	if err := doc.GenMarkdownTree(cmd.NewRootCmd(), docPath); err != nil {
		log.Fatal(err)
	}
}
