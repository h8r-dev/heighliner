package dagger

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-getter/v2"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

const installScriptURL = "https://dl.dagger.io/dagger/install.sh"

// Client maintains dagger binary and executes dagger commands.
type Client struct {
	// Path to the executable binary file of dagger.
	Binary string

	LogFormat string
	LogLevel  string

	genericclioptions.IOStreams
}

// ActionOptions controls the behavior of dagger action.
type ActionOptions struct {
	// Name of the action to be executed.
	Name string

	// Consider the following example:
	//
	// .
	// └── example
	// 	├── cue.mod
	// 	│   ├── module.cue
	// 	│   ├── pkg
	// 	│   └── usr
	// 	└── plans
	// 		└── plan.cue
	//
	// 'Dir' should be ./ and 'Plan' should be ./plans
	// or ./plans/plan.cue (You can specify plan to
	// both a directory or a file).

	// Dir is the path to your cue module (the parent dir of
	// the 'cue.mod' dir that contains 'module.cue' file).
	Dir string
	// Relative path from `dir` to your plan, which is
	// expected to begin with `.` (default ".").
	Plan string
	// Disable caching when `NoCache` is set to `true`.
	NoCache bool
}

// GetPath returns the path to the dagger binary.
func GetPath() string {
	return filepath.Join(state.GetHln(), "bin", "dagger")
}

// NewClient creates a customized dagger client and returns it
func NewClient(logFormat, logLevel string, streams genericclioptions.IOStreams) (*Client, error) {
	binary := GetPath()
	return &Client{
		Binary:    binary,
		LogFormat: logFormat,
		LogLevel:  logLevel,
		IOStreams: streams,
	}, nil
}

// NewDefaultClient creates a default dagger client and returns it.
func NewDefaultClient(streams genericclioptions.IOStreams) (*Client, error) {
	binary := GetPath()
	return &Client{
		Binary:    binary,
		LogFormat: "plain",
		LogLevel:  "info",
		IOStreams: streams,
	}, nil
}

// NewActionOptions creates and returns a ActionOptions struct.
func NewActionOptions(name, dir, plan string, noCache bool) *ActionOptions {
	return &ActionOptions{
		Name:    name,
		Dir:     dir,
		Plan:    plan,
		NoCache: noCache,
	}
}

// Do executes a dagger do command.
func (c *Client) Do(o *ActionOptions) error {
	if o.Dir != "" {
		err := os.Chdir(o.Dir)
		if err != nil {
			return err
		}
	}
	args := []string{
		"--log-format", c.LogFormat,
		"--log-level", c.LogLevel,
		"do", o.Name,
		"--plan", o.Plan,
	}
	if o.NoCache {
		args = append(args, "--no-cache")
	}
	return util.Exec(c.IOStreams, c.Binary, args...)
}

// Check installs or upgrades the dagger binary.
func (c *Client) Check() error {
	// Check if dagger binary exist.
	if _, err := os.Stat(c.Binary); errors.Is(err, os.ErrNotExist) {
		return c.installOrUpgrade()
	}
	// Check if the version of dagger is the latest.
	buf := &bytes.Buffer{}
	omw := io.MultiWriter(buf, c.IOStreams.Out)
	emw := io.MultiWriter(buf, c.IOStreams.ErrOut)
	err := util.Exec(genericclioptions.IOStreams{
		In:     buf,
		Out:    omw,
		ErrOut: emw,
	}, GetPath(), "version", "--check")
	if err != nil {
		return err
	}
	msg := buf.String()
	if !strings.Contains(msg, "up to date") {
		return c.installOrUpgrade()
	}
	return nil
}

// installOrUpgrade run the dagger install.sh script.
func (c *Client) installOrUpgrade() error {
	dir := filepath.Dir(filepath.Dir(c.Binary))
	err := os.MkdirAll(dir, 0755)
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
