package state

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var (
	ErrNoStack = fmt.Errorf("no stack in datastore")
)

type DataStore struct {
	Path string `json:"path"`
}

// Make makes hln dir in dst
func Make(dst string) (*DataStore, error) {
	dir := path.Join(dst, "hln")
	if dst == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, errors.New("failed to get current working dir")
		}
		dir = filepath.Join(cwd, "hln")
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to make dir (%s): %w", dir, err)
	}
	var ds = &DataStore{
		Path: dir,
	}
	return ds, nil
}

// Find the stack in the current datastore dir
func (ds *DataStore) Find() (*Stack, error) {
	dir, err := ioutil.ReadDir(ds.Path)
	if err != nil {
		return nil, err
	}
	s := &Stack{
		Path: "",
	}
	for _, fi := range dir {
		if fi.IsDir() {
			s.Name = fi.Name()
			s.Path = path.Join(ds.Path)
			break
		}
	}
	if s.Path == "" {
		return nil, ErrNoStack
	}
	return s, nil
}
