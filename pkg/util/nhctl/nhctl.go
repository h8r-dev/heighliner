package nhctl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-getter/v2"

	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

const version = "v0.6.16"

// GetPath returns the path to the nhctl binary.
func GetPath() string {
	return filepath.Join(state.GetHln(), "bin", getName())
}

func getName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	name := fmt.Sprintf("nhctl-%s-%s", os, arch)
	if os == "windows" {
		name += ".exe"
	}
	return name
}

// Check install nhctl binary if necessary.
func Check() error {
	var err error
	hln := state.GetHln()
	nhctl := GetPath()

	_, err = os.Stat(nhctl)
	if err != nil {
		return update(hln)
	}
	return nil
}

func update(dir string) error {
	src := fmt.Sprintf("https://github.com/nocalhost/nocalhost/releases/download/%s/%s", version, getName())
	req := &getter.Request{
		Src: src,
		Dst: filepath.Join(dir, "bin"),
	}
	err := util.GetWithTracker(req)
	if err != nil {
		return err
	}
	err = os.Chmod(filepath.Join(dir, "bin", getName()), 0700)
	if err != nil {
		return err
	}
	return nil
}
