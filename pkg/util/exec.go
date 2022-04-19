package util

import (
	"os/exec"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// Exec executes the command and prints the output into current terminal
func Exec(streams genericclioptions.IOStreams, name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Stdin = streams.In
	cmd.Stdout = streams.Out
	cmd.Stderr = streams.ErrOut

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
