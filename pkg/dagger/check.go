package dagger

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/go-getter/v2"
	gover "github.com/hashicorp/go-version"
	"go.uber.org/zap"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/version"
)

const installScriptURL = "https://dl.dagger.io/dagger/install.sh"

// Check checks if the version of dagger binary is available.
func (c *Client) Check() error {
	// Check if dagger binary exist.
	if _, err := os.Stat(c.Binary); errors.Is(err, os.ErrNotExist) {
		return errors.New("no dagger binary file found")
	}
	// Check if the version of dagger is the latest.
	rex := regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)
	buf := &bytes.Buffer{}
	err := util.Exec(genericclioptions.IOStreams{
		In:     buf,
		Out:    buf,
		ErrOut: buf,
	}, GetPath(), "version")
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
	err := os.Setenv("DAGGER_VERSION", version.DaggerDefault)
	if err != nil {
		return err
	}
	dir := filepath.Dir(filepath.Dir(c.Binary))
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	installScript := filepath.Join(dir, "dagger", "install.sh")
	_, err = os.Stat(installScript)
	if err == nil {
		err = os.Remove(installScript)
		if err != nil {
			return err
		}
	}
	req := &getter.Request{
		Src: installScriptURL,
		Dst: filepath.Dir(installScript),
	}
	err = util.GetWithTracker(req)
	if err != nil {
		return err
	}
	err = util.Exec(c.IOStreams, "/bin/sh", installScript)
	if err != nil {
		return err
	}
	return nil
}
