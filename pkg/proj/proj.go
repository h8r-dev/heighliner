package proj

import (
	"errors"
	"path"

	"github.com/otiai10/copy"

	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/h8r-dev/heighliner/pkg/state"
)

// Proj holds information about it's stack and path.
type Proj struct {
	Stack *stack.Stack
	Home  state.Storage
}

// New creates a project and returns it.
func New(s *stack.Stack, h state.Storage) *Proj {
	return &Proj{
		Stack: s,
		Home:  h,
	}
}

// Init detects if there is a project already.
// If there isn's any then creates one and initializes it.
func (p *Proj) Init() error {
	s := p.Stack
	if err := p.Home.Detect(); err == nil {
		return errors.New("already have a project")
	}
	if err := state.CleanCache(s); err != nil {
		return err
	}
	if err := s.Pull(state.Cache); err != nil {
		return err
	}
	if err := p.Home.Init(); err != nil {
		return err
	}
	if err := copy.Copy(path.Join(state.Cache, p.Stack.Name), p.Home.GetPath()); err != nil {
		return err
	}
	return p.Home.Detect()
}

// Detect calls the Detect method of project's Home.
func (p *Proj) Detect() error {
	return p.Home.Detect()
}

// Drop cleans the current project.
func (p *Proj) Drop() error {
	return p.Home.Clean()
}
