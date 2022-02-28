package stack

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/h8r-dev/heighliner/pkg/util/compress"
	"gopkg.in/yaml.v3"
)

type Stack struct {
	Name        string `json:"name" yaml:"name"`
	Path        string `json:"path"`
	Url         string `json:"url"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

var (
	Sample = &Stack{
		Name:        "sample",
		Url:         "https://stack.h8r.io/sample-latest.tar.gz",
		Description: "This is an example stack",
	}
	GoGinStack = &Stack{
		Name:        "go-gin-stack",
		Url:         "https://stack.h8r.io/go-gin-stack-latest.tar.gz",
		Description: "This is an go-gin stack",
	}
)

var Stacks = map[string]*Stack{
	"sample":       Sample,
	"go-gin-stack": GoGinStack,
}

// New creates a Stack struct and a dir to store it's files
func New(name, dst, src string) (*Stack, error) {
	dir := filepath.Join(dst, name)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to make dir %s in %s", name, dst)
	}
	s := &Stack{
		Name: name,
		Path: dst,
		Url:  src,
	}
	return s, nil
}

// Download downloads the stack form it's Url field
func (s *Stack) Download() error {
	fp := filepath.Join(s.Path, s.Name+".tar.gz")
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	rsp, err := http.Get(s.Url)
	if err != nil {
		return fmt.Errorf("failed to download stack from %s: %w", s.Url, err)
	}
	defer rsp.Body.Close()

	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Decompress decompresses the raw .tar.gz package of a stack
func (s *Stack) Decompress() error {
	src := filepath.Join(s.Path, s.Name+".tar.gz")

	err := compress.Decompress(src, s.Path)
	if err != nil {
		return fmt.Errorf("failed to decompress stack %s: %w", s.Name, err)
	}

	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}

// Load() loads values from the metadata.yaml file
func (s *Stack) Load() error {
	metadata := path.Join(s.Path, "metadata.yaml")
	file, err := os.Open(metadata)
	if err != nil {
		return err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, s)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data from %s: %w", metadata, err)
	}
	return nil
}
