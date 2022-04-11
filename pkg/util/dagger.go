package util

import (
	"os"
	"path"

	"github.com/hashicorp/go-getter/v2"

	"github.com/h8r-dev/heighliner/pkg/state"
)

var (
	Dagger           = path.Join(state.GetHln(), "bin", "dagger")
	installDaggerSrc = "https://dl.dagger.io/dagger/install.sh"
)

// CheckDagger install dagger binary if necessary.
func CheckDagger() error {
	var (
		err        error
		hln        = state.GetHln()
		daggerFile = path.Join(hln, "bin", "dagger")
	)

	_, err = os.Stat(daggerFile)
	if err != nil {
		return updateDagger(hln)
	}
	return nil
}

func updateDagger(dir string) error {
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
	err = GetWithTracker(req)
	if err != nil {
		return err
	}
	err = Exec("/bin/sh", path.Join("dagger", "install.sh"))
	if err != nil {
		return err
	}
	return nil
}
