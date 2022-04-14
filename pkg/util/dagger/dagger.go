package dagger

import (
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/go-getter/v2"

	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

const installDaggerSrc = "https://dl.dagger.io/dagger/install.sh"

// GetPath represents the path to the dagger binary.
func GetPath() string {
	return filepath.Join(state.GetHln(), "bin", "dagger")
}

// Check install dagger binary if necessary.
func Check() error {
	var (
		err        error
		hln        = state.GetHln()
		daggerFile = path.Join(hln, "bin", "dagger")
	)

	_, err = os.Stat(daggerFile)
	if err != nil {
		return update(hln)
	}
	return nil
}

func update(dir string) error {
	var err error
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	req := &getter.Request{
		Src: installDaggerSrc,
		Dst: path.Join(dir, "dagger"),
	}
	err = util.GetWithTracker(req)
	if err != nil {
		return err
	}
	err = util.Exec("/bin/sh", path.Join("dagger", "install.sh"))
	if err != nil {
		return err
	}
	return nil
}
