package state

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/go-getter/v2"
	"github.com/otiai10/copy"
)

// Stack is a CloudNative app template
type Stack struct {
	Name        string `json:"name" yaml:"name"`
	Path        string `json:"path"`
	URL         string `json:"url"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

var (
	// ErrStackNotExist means heighliner can't find the stack in localstorage
	ErrStackNotExist = errors.New("target stack doesn't exist")
)

var (
	// SampleStack is a demo stack that echoes your input
	SampleStack = &Stack{
		Name:        "sample",
		URL:         "https://stack.h8r.io/sample-latest.tar.gz",
		Description: "Sample is a light-weight stack mainly used for test",
		Version:     "1.0.0",
	}
	// GoGinStack is a go microservice architecture app
	GoGinStack = &Stack{
		Name:        "go-gin-stack",
		URL:         "https://stack.h8r.io/go-gin-stack-latest.tar.gz",
		Description: "Go-gin-stack helps you configure many cloud native components including prometheus, grafana, nocalhost, etc.",
		Version:     "1.0.0",
	}
)

// Stacks stores all stacks that currently usable
var Stacks = map[string]*Stack{
	"sample":       SampleStack,
	"go-gin-stack": GoGinStack,
}

// CleanStacks cleans all cached stacks
func CleanStacks() error {
	err := initHeighlinerCache()
	if err != nil {
		return err
	}
	dir := path.Join(HeighlinerCacheHome, "repository")
	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

// NewStack returns a Stack struct
func NewStack(name string) (*Stack, error) {
	err := initHeighlinerCache()
	if err != nil {
		return nil, err
	}
	return &Stack{
		Name: name,
		Path: path.Join(HeighlinerCacheHome, "repository"),
	}, nil
}

// Pull downloads and decompresses a stack
func (s *Stack) Pull(url string) error {
	req := &getter.Request{
		Src: url,
		Dst: s.Path,
	}
	err := getWithTracker(req)
	if err != nil {
		return fmt.Errorf("failed to download stack: %w", err)
	}
	return nil
}

// Check checks the status of target stack
func (s *Stack) Check() error {
	dir := path.Join(s.Path, s.Name)

	_, err := os.Stat(dir)
	if err != nil {
		return ErrStackNotExist
	}

	return nil
}

// Copy copies the stack into dest
func (s *Stack) Copy(dest string) error {
	src := path.Join(s.Path, s.Name)
	err := copy.Copy(src, dest)
	if err != nil {
		return fmt.Errorf("failed to copy stack %s: %w", s.Name, err)
	}
	return nil
}
