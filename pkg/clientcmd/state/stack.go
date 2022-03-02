package state

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/h8r-dev/heighliner/pkg/clientcmd/util"
	"github.com/otiai10/copy"
)

type Stack struct {
	Name        string `json:"name" yaml:"name"`
	Path        string `json:"path"`
	Url         string `json:"url"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

var (
	ErrStackNotExist = errors.New("target stack doesn't exist")
)

var (
	SampleStack = &Stack{
		Name:        "sample",
		Url:         "https://stack.h8r.io/sample-latest.tar.gz",
		Description: "This is an example stack",
		Version:     "1.0.0",
	}
	GoGinStack = &Stack{
		Name:        "go-gin-stack",
		Url:         "https://stack.h8r.io/go-gin-stack-latest.tar.gz",
		Description: "This is an go-gin stack",
		Version:     "1.0.0",
	}
)

var Stacks = map[string]*Stack{
	"sample":       SampleStack,
	"go-gin-stack": GoGinStack,
}

// New creates a Stack struct and a dir to store it's files
func NewStack(name string) *Stack {
	return &Stack{
		Name: name,
	}
}

// Pull downloads and decompresses a stack
func (s *Stack) Pull(url string) error {
	s.Url = url

	uhd, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}

	dir := path.Join(uhd, ".hln")

	err = os.MkdirAll(dir, 0755)
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

//Check checks the status of target stack
func (s *Stack) Check() error {
	uhd, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}

	dir := path.Join(uhd, ".hln", s.Name)

	_, err = os.Stat(dir)
	if err != nil {
		return ErrStackNotExist
	}

	s.Path = path.Join(uhd, ".hln")
	return nil
}

// Cpoy copies stack into dest
func (s *Stack) Copy(dest string) error {
	src := path.Join(s.Path, s.Name)
	err := copy.Copy(src, dest)
	if err != nil {
		return fmt.Errorf("failed to copy stack %s: %w", s.Name, err)
	}
	return nil
}

func (s *Stack) download() error {
	fp := filepath.Join(s.Path, s.Name+".tar.gz")
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	rsp, err := http.Get(s.Url)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *Stack) decompress() error {
	src := filepath.Join(s.Path, s.Name+".tar.gz")

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
