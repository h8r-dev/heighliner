package stack

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-getter/v2"
	"github.com/otiai10/copy"

	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

// Stack is a CloudNative app template
type Stack struct {
	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url" yaml:"url"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

var (
	// ErrNoSuchStack means heighliner can't find the stack in localstorage
	ErrNoSuchStack = errors.New("target stack doesn't exist")
)

// Stacks stores all stacks that currently usable
var Stacks = map[string]struct{}{
	"sample":       {},
	"go-gin-stack": {},
	"gin-vue":      {},
}

// New returns a Stack struct
func New(name string) (*Stack, error) {
	const defaultVersion = "latest"

	// Check if specified stack exists or not
	_, ok := Stacks[name]
	if !ok {
		return nil, ErrNoSuchStack
	}

	version := defaultVersion
	url := fmt.Sprintf("https://stack.h8r.io/%s-%s.tar.gz", name, version)
	s := &Stack{
		Name:    name,
		URL:     url,
		Version: version,
	}
	return s, nil
}

// Pull downloads and extracts the stack
func (s *Stack) Pull() error {
	req := &getter.Request{
		Src: s.URL,
		Dst: state.Cache,
	}
	err := util.GetWithTracker(req)
	if err != nil {
		state.CleanCache()
		return fmt.Errorf("failed to pull stack: %w", err)
	}
	return nil
}

// Copy the stack into dst dir
func (s *Stack) Copy(src, dst string) error {
	return copy.Copy(src, dst)
}
