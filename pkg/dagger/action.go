package dagger

import (
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

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
	if err := util.Exec(genericclioptions.NewTestIOStreamsDiscard(), c.Binary, "project", "init"); err != nil {
		return err
	}
	if err := util.Exec(c.IOStreams, c.Binary, "project", "update"); err != nil {
		return err
	}
	// For convenience that user might forget to set KUBECONFIG env, we will still set it which our stacks depends on.
	if os.Getenv("KUBECONFIG") == "" {
		_ = os.Setenv("KUBECONFIG", k8sutil.GetKubeConfigPath())
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
