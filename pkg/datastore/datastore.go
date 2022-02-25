package datastore

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/h8r-dev/heighliner/pkg/stack"
)

type DataStore struct {
	Path string `json:"path"`
}

// Init creates datastore in the currernt workdir
func Init() (*DataStore, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("failed to get current working dir")
	}
	dir := filepath.Join(cwd, ".hln")
	// Make .hln Dir
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to make dir (%s): %w", dir, err)
	}
	ds := &DataStore{
		Path: "",
	}
	ds.Path = dir
	return ds, nil
}

// Stat shows the status of datastore
func Stat() (*DataStore, error) {

	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("failed to get current working dir")
	}
	dir := filepath.Join(cwd, ".hln")

	_, err = os.Stat(dir)
	if err != nil {
		return nil, errors.New("please init firstly")
	}
	var ds = &DataStore{
		Path: "",
	}
	ds.Path = dir
	return ds, nil
}

// Find the stack in the current datastore dir
func (ds *DataStore) Find() (*stack.Stack, error) {
	dir, err := ioutil.ReadDir(ds.Path)
	if err != nil {
		return nil, err
	}
	s := &stack.Stack{
		Path: "",
	}
	for _, fi := range dir {
		if fi.IsDir() && fi.Name() != "env" {
			s.Path = path.Join(ds.Path, fi.Name())
		}
	}
	if s.Path == "" {
		return nil, errors.New("can not find a stack in current space")
	}
	return s, nil
}
