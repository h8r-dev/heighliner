package state

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type State struct {
	Env string `json:"env"`
}

// InitEnv initializes the local state, e.g. environment data, cue mods.
func InitEnv(name string) (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("failed to get current working dir")
	}
	// prepare env dir
	dir := filepath.Join(cwd, ".hln", "env")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to make dir (%s): %w", dir, err)
	}

	var s *State
	fn := filepath.Join(dir, name)
	_, err = os.Stat("/path/to/whatever")
	switch {
	case errors.Is(err, os.ErrNotExist):
		// create a new env state file
		s = &State{
			Env: name,
		}
		b, err := yaml.Marshal(s)
		if err != nil {
			return nil, err
		}
		err = ioutil.WriteFile(fn, b, 0600)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, fmt.Errorf("stat file (%s): %w", fn, err)
	default:
		// load existing env state
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(b, s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}
