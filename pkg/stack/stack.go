package stack

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-getter/v2"
	"github.com/otiai10/copy"

	"github.com/h8r-dev/heighliner/pkg/util"
)

// Stack is a CloudNative app template
type Stack struct {
	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

var (
	// ErrNoSuchStack means heighliner can't find the stack in localstorage
	ErrNoSuchStack = errors.New("target stack doesn't exist")
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
		Description: "go-gin-stack helps you configure many cloud native components including prometheus, grafana, nocalhost, etc.",
		Version:     "1.0.0",
	}
	// GinVueStack is new version of go-gin-stack
	GinVueStack = &Stack{
		Name:        "gin-vue",
		URL:         "https://stack.h8r.io/gin-vue-latest.tar.gz",
		Description: "gin-vue is a new version of go-gin-stack",
		Version:     "1.0.0",
	}
)

// Stacks stores all stacks that currently usable
var Stacks = map[string]*Stack{
	"sample":       SampleStack,
	"go-gin-stack": GoGinStack,
	"gin-vue":      GinVueStack,
}

// New returns a Stack struct
func New(name string) (*Stack, error) {
	// Check if specified stack exist or not
	val, ok := Stacks[name]
	if !ok {
		return nil, ErrNoSuchStack
	}
	return val, nil
}

// Pull downloads and decompresses a stack
func (s *Stack) Pull(dst string) error {
	req := &getter.Request{
		Src: s.URL,
		Dst: dst,
	}
	err := util.GetWithTracker(req)
	if err != nil {
		return fmt.Errorf("failed to pull stack: %w", err)
	}
	return nil
}

// Copy the stack into dst dir
func (s *Stack) Copy(src, dst string) error {
	return copy.Copy(src, dst)
}
