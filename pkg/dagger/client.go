package dagger

import (
	"path/filepath"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/state"
)

// Client maintains dagger binary and executes dagger commands.
type Client struct {
	// Path to the executable binary file of dagger.
	Binary string

	LogFormat string
	LogLevel  string

	genericclioptions.IOStreams
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
