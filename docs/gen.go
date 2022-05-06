package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra/doc"

	"github.com/h8r-dev/heighliner/cmd/hln/cmd"
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
	if err := gennerateDocs(docPath); err != nil {
		panic(err)
	}
	if err := quoteFiglet(docPath); err != nil {
		panic(err)
	}
}

func gennerateDocs(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	return doc.GenMarkdownTree(cmd.NewRootCmd(), path)
}

func quoteFiglet(path string) error {
	src := filepath.Join(path, "hln.md")
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		if line == "### Synopsis" {
			lines[i+2] = "```"
		}
		if line == "### Options" {
			lines[i-2] = "```"
		}
	}
	output := strings.Join(lines, "\n")
	return os.WriteFile(src, []byte(output), 0644)
}
