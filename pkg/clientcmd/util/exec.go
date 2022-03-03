package util

import (
	"fmt"
	"os"
	"os/exec"
)

// Exec executes the command and prints the output into current terminal
func Exec(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run %s: %w", name, err)
	}
	return nil
}
