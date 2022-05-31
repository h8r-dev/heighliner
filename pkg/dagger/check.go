package dagger

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
	"github.com/h8r-dev/heighliner/pkg/util/ziputil"
	"github.com/h8r-dev/heighliner/pkg/version"
)

const (
	installScriptURL = "https://dl.dagger.io/dagger/install.sh"
	windowsBase      = "https://dagger-io.s3.amazonaws.com"
)

// Check checks if the version of dagger binary is available.
func (c *Client) Check() error {
	lg := logger.New(c.IOStreams)
	// Check if dagger binary exist.
	if _, err := os.Stat(c.Binary); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("no dagger binary file found in %s", c.Binary)
	}
	// Check if the version of dagger is the available.
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
	constraints, err := gover.NewConstraint(version.DaggerConstraint)
	if err != nil {
		return err
	}
	if !constraints.Check(ver) {
		return fmt.Errorf("current dagger version: %s, expect %s",
			ver, version.DaggerConstraint)
	}
	lg.Info(fmt.Sprintf("dagger version %s", ver.Original()))
	return nil
}

// CheckAndInstall installs dagger if necessary.
func (c *Client) CheckAndInstall() error {
	lg := logger.New(c.IOStreams)
	if err := c.Check(); err != nil {
		lg.Info("downloading dagger...", zap.NamedError("info", err))
		return c.install()
	}
	return nil
}

// install runs the dagger install.sh script.
func (c *Client) install() error {
	if runtime.GOOS == "windows" {
		return c.installForWindows()
	}
	if err := os.Setenv("DAGGER_VERSION", version.DaggerDefault); err != nil {
		return err
	}
	dst := filepath.Dir(filepath.Dir(c.Binary))
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	if err := os.Chdir(dst); err != nil {
		return err
	}
	src := installScriptURL
	shName := "installDagger.sh"
	if err := getter.Get(c.Out, getter.NewRequest(src, dst, shName)); err != nil {
		return err
	}
	shFile := filepath.Join(dst, shName)
	if err := util.Exec(c.IOStreams, "/bin/sh", shFile); err != nil {
		return err
	}
	return os.Remove(shFile)
}

func (c Client) installForWindows() error {
	fileName := "dagger_v" + version.DaggerDefault + "_windows_amd64"
	src := fmt.Sprintf("%s/dagger/releases/%s/%s.zip", windowsBase, version.DaggerDefault, fileName)
	dst := filepath.Dir(c.Binary)
	zipName := "dagger.zip"
	if err := getter.Get(c.Out, getter.NewRequest(src, dst, zipName)); err != nil {
		return err
	}
	zipFile := filepath.Join(dst, zipName)
	if err := ziputil.Extract(dst, zipFile); err != nil {
		return err
	}
	return os.Remove(zipFile)
}
