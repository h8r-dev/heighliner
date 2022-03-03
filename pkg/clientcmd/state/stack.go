package state

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/otiai10/copy"
	"github.com/rs/zerolog/log"

	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
)

// Stack is a CloudNative app template
type Stack struct {
	Name        string `json:"name" yaml:"name"`               // The Name of the stack
	Path        string `json:"path"`                           // Path to localstorage
	URL         string `json:"url"`                            // URL to pull the stack
	Version     string `json:"version" yaml:"version"`         // Version
	Description string `json:"description" yaml:"description"` // A short description
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
		Description: "This is an example stack",
		Version:     "1.0.0",
	}
	// GoGinStack is a go microservice architecture app
	GoGinStack = &Stack{
		Name:        "go-gin-stack",
		URL:         "https://stack.h8r.io/go-gin-stack-latest.tar.gz",
		Description: "This is an go-gin stack",
		Version:     "1.0.0",
	}
)

// Stacks stores all stacks that currently usable
var Stacks = map[string]*Stack{
	"sample":       SampleStack,
	"go-gin-stack": GoGinStack,
}

// NewStack returns a Stack struct
func NewStack(name string) *Stack {
	return &Stack{
		Name: name,
	}
}

// Pull downloads and decompresses a stack
func (s *Stack) Pull(url string) error {
	s.URL = url

	uhd, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}

	dir := path.Join(uhd, ".hln")

	err = os.MkdirAll(dir, 0750)
	if err != nil {
		return fmt.Errorf("failed to initialize local storage: %w", err)
	}

	s.Path = dir

	err = s.download()
	if err != nil {
		return fmt.Errorf("failed to download stack %s: %w", s.Name, err)
	}

	err = s.decompress()
	if err != nil {
		return fmt.Errorf("failed to decompress stack %s: %w", s.Name, err)
	}

	return nil
}

// Check checks the status of target stack
func (s *Stack) Check() error {
	uhd, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}

	s.Path = path.Join(uhd, ".hln")

	dir := path.Join(uhd, ".hln", s.Name)

	_, err = os.Stat(dir)
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

func (s *Stack) download() error {
	fp := path.Join(s.Path, "temp.tar.gz")
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()

	rsp, err := http.Get(s.URL)
	if err != nil {
		return err
	}
	defer func() {
		if err := rsp.Body.Close(); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()

	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *Stack) decompress() error {
	src := filepath.Join(s.Path, "temp.tar.gz")

	err := util.Decompress(src, s.Path)
	if err != nil {
		return err
	}

	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}
