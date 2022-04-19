package project

// Projects only exist internally.

import (
	"os"

	"github.com/otiai10/copy"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/util"
)

// Project is a dir where dagger plan is executed.
type Project struct {
	Home string
	Src  string
}

// New creates a Project object and returns it.
func New(src, home string) *Project {
	return &Project{
		Home: home,
		Src:  src,
	}
}

// Init initializes the project.
func (p *Project) Init() error {
	var err error

	p.clean()
	err = copy.Copy(p.Src, p.Home)
	if err != nil {
		return err
	}

	err = p.init()
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) init() error {
	var err error

	err = os.Chdir(p.Home)
	if err != nil {
		return err
	}

	// Initialize & update dagger project.
	// TODO move this to dagger client package.
	err = util.Exec(genericclioptions.NewTestIOStreamsDiscard(), dagger.GetPath(), "project", "init")
	if err != nil {
		return err
	}
	err = util.Exec(genericclioptions.NewTestIOStreamsDiscard(), dagger.GetPath(), "project", "update")
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) clean() {
	if err := os.RemoveAll(p.Home); err != nil {
		panic(err)
	}
}
