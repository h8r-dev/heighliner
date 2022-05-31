package nhctl

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	gover "github.com/hashicorp/go-version"
	"go.uber.org/zap"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/util/getter"
	"github.com/h8r-dev/heighliner/pkg/version"
)

// Check checks if the nhctl version is available.
func (c *Client) Check() error {
	lg := logger.New(c.IOStreams)
	// Check if nhctl binary exist.
	if _, err := os.Stat(c.Binary); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("no nhctl binary file found in %s", c.Binary)
	}
	// Check if the version of nhctl is the available.
	rex := regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)
	buf := &bytes.Buffer{}
	err := util.Exec(genericclioptions.IOStreams{
		In:     buf,
		Out:    buf,
		ErrOut: buf,
	}, GetBin(), "version")
	if err != nil {
		return err
	}
	msg := buf.String()
	ver, err := gover.NewSemver(rex.FindString(msg))
	if err != nil {
		return err
	}
	constraints, err := gover.NewConstraint(version.NhctlConstraint)
	if err != nil {
		return err
	}
	if !constraints.Check(ver) {
		return fmt.Errorf("current nhctl version: %s, expect %s",
			ver, version.NhctlConstraint)
	}
	lg.Info(fmt.Sprintf("nhctl version %s", ver.Original()))
	return nil
}

// CheckAndInstall will install nhctl if necessary.
func (c *Client) CheckAndInstall() error {
	lg := logger.New(c.IOStreams)
	if err := c.Check(); err != nil {
		lg.Info("downloading nhctl...", zap.NamedError("info", err))
		return c.install()
	}
	return nil
}

func (c *Client) install() error {
	src := fmt.Sprintf(
		"https://github.com/nocalhost/nocalhost/releases/download/v%s/%s",
		version.NhctlDefault, c.getName())
	dst := filepath.Dir(c.Binary)
	if err := getter.Get(os.Stdout, getter.NewRequest(src, dst, c.getName())); err != nil {
		return err
	}
	if err := os.Rename(filepath.Join(filepath.Dir(c.Binary), c.getName()),
		c.Binary); err != nil {
		return err
	}
	if err := os.Chmod(c.Binary, 0700); err != nil {
		return err
	}
	return nil
}

func (c *Client) getName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	name := fmt.Sprintf("nhctl-%s-%s", os, arch)
	if os == "windows" {
		name += ".exe"
	}
	return name
}
