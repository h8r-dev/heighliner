package nhctl

import (
	"path/filepath"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/state"
)

// Client interactive with nhctl
type Client struct {
	// Path to the executable binary file of nhctl.
	Binary string

	genericclioptions.IOStreams
}

// GetBin returns the path to the nhctl binary.
func GetBin() string {
	return filepath.Join(state.GetHln(), "bin", "nhctl")
}

// NewDefaultClient creates a default nhctl client and returns it.
func NewDefaultClient(streams genericclioptions.IOStreams) (*Client, error) {
	return &Client{
		Binary:    GetBin(),
		IOStreams: streams,
	}, nil
}
