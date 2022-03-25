package state

import (
	"errors"
	"fmt"
	"os"
	"path"
)

// Temp is the place to store and execute our project
type Temp struct {
	Path string
}

// NewTemp creates and returns a Temp object
func NewTemp() *Temp {
	t := &Temp{}
	t.Path = os.Getenv("HLN_TEMP_HOME")
	if t.Path == "" {
		tempDir := os.TempDir()
		t.Path = path.Join(tempDir, "heighliner")
	}
	return t
}

// Init initializes the temporary localstorage
func (t *Temp) Init() error {
	if err := t.Detect(); err == nil {
		return errors.New("tmp dir already exists")
	}
	if err := os.MkdirAll(t.Path, 0755); err != nil {
		return fmt.Errorf("failed to create dir %s: %w", t.Path, err)
	}
	return nil
}

// GetPath fetches the path of localstorage
func (t *Temp) GetPath() string {
	return t.Path
}

// Clean deletes the temporary localstorage
func (t *Temp) Clean() error {
	return os.RemoveAll(t.Path)
}

// Detect detects if the temporary exists.
// If it does exist, then enter it.
func (t *Temp) Detect() error {
	_, err := os.Stat(t.Path)
	if err != nil {
		return err
	}
	return os.Chdir(t.GetPath())
}
