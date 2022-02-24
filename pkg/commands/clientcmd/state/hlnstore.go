package state

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type HlnStore struct {
	Path   string   `json:"path"`
	Stacks []*Stack `json:"stacks"`
}

// InitHlnStore initializes the .hln dir to keep stacks and other things
func InitHlnStore() (*HlnStore, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.New("failed to get user home dir")
	}
	// Create the .hln dir
	dir := filepath.Join(userHomeDir, ".hln")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to make dir (%s): %w", dir, err)
	}
	hs := &HlnStore{
		Path: dir,
	}
	return hs, nil
}

// Add a new stack into HlnStore
func (hs *HlnStore) NewStack(name, url string) (*Stack, error) {
	dir := filepath.Join(hs.Path, "stacks", name)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	s := &Stack{
		Name: name,
		Url:  url,
		Path: dir,
	}
	hs.Stacks = append(hs.Stacks, s)
	return s, nil
}
