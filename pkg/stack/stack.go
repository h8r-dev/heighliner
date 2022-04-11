package stack

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/go-getter/v2"

	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

// Stack is a CloudNative application template.
type Stack struct {
	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url" yaml:"url"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

var (
	// ErrNoSuchStack means heighliner can't find the stack in localstorage.
	ErrNoSuchStack = errors.New("target stack doesn't exist")
)

// Stacks stores all available stacks.
var Stacks = map[string]struct{}{
	"gin-vue":  {},
	"gin-next": {},
}

// New returns a Stack object.
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

// Update upgrades the stack if necessary.
func (s *Stack) Update() error {
	ok := s.check()
	if !ok {
		s.clean()
		if err := s.pull(); err != nil {
			s.clean()
			return err
		}
	}
	return nil
}

// check checks if the stack is up to date.
func (s *Stack) check() bool {
	return false
}

func (s *Stack) pull() error {
	req := &getter.Request{
		Src: s.URL,
		Dst: state.GetCache(),
	}
	err := util.GetWithTracker(req)
	if err != nil {
		return err
	}
	return nil
}

func (s *Stack) clean() {
	if err := os.RemoveAll(path.Join(state.GetCache(), s.Name)); err != nil {
		panic(err)
	}
}
