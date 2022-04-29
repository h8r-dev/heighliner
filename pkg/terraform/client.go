package terraform

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/hashicorp/go-getter/v2"
	gover "github.com/hashicorp/go-version"
	"go.uber.org/zap"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/version"
)

// Client interactive with terraform
type Client struct {
	// Path to the executable binary file of terraform.
	Binary string

	genericclioptions.IOStreams
}

// GetBin returns the path to the terraform binary.
func GetBin() string {
	return filepath.Join(state.GetHln(), "bin", "terraform")
}

// NewDefaultClient creates a default terraform client.
func NewDefaultClient(streams genericclioptions.IOStreams) (*Client, error) {
	return &Client{
		Binary:    GetBin(),
		IOStreams: streams,
	}, nil
}

// Check checks if the version of terraform binary is available.
func (c *Client) Check() error {
	lg := logger.New(c.IOStreams)
	// Check if terraform binary exist.
	if _, err := os.Stat(c.Binary); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("no terraform binary file found in %s", c.Binary)
	}
	// Check if the version of terraform is the available.
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
	constraints, err := gover.NewConstraint(version.TerraformConstraint)
	if err != nil {
		return err
	}
	if !constraints.Check(ver) {
		return fmt.Errorf("current terraform version: %s, expect %s",
			ver, version.TerraformConstraint)
	}
	lg.Info(fmt.Sprintf("terraform version %s", ver.Original()))
	return nil
}

// CheckAndInstall will install terraform if necessary.
func (c *Client) CheckAndInstall() error {
	lg := logger.New(c.IOStreams)
	if err := c.Check(); err != nil {
		lg.Info("downloading terraform...", zap.NamedError("info", err))
		return c.install()
	}
	return nil
}

func (c *Client) install() error {
	src := fmt.Sprintf(
		"https://dl.h8r.io/terraform/terraform_%s_%s_%s.zip",
		version.TerraformDefault, runtime.GOOS, runtime.GOARCH)
	req := &getter.Request{
		Src: src,
		Dst: filepath.Dir(c.Binary),
	}
	err := util.GetWithTracker(req)
	if err != nil {
		return err
	}
	return nil
}
