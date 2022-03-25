package state

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
)

var (
	// Temp is the place to store and execute our project
	Temp string
)

func init() {
	Temp = os.Getenv("HLN_TEMP_HOME")
	if Temp == "" {
		tempDir := os.TempDir()
		Temp = path.Join(tempDir, "heighliner")
	}
}

// InitTemp creates and returns a Temp object
func InitTemp() error {
	if err := EnterTemp(); err == nil {
		return errors.New("failed to initialize temp dir")
	}
	if err := os.MkdirAll(Temp, 0755); err != nil {
		return fmt.Errorf("failed to create dir %s: %w", Temp, err)
	}
	return nil
}

// EnterTemp tries to enter the Temp dir
// and returns an error if it doesn't exist.
func EnterTemp() error {
	_, err := os.Stat(Temp)
	if err != nil {
		return err
	}
	return os.Chdir(Temp)
}

// CleanTemp removes all things in Temp dir
func CleanTemp() error {
	return os.RemoveAll(Temp)
}
