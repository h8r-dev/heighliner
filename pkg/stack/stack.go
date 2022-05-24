package stack

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter/v2"
	"gopkg.in/yaml.v2"

	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

// Stack is a CloudNative application template.
type Stack struct {
	Path string

	// TODO Should read from stack metadata
	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url" yaml:"url"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

// StacksIndexURL to get index
const StacksIndexURL = "https://stack.h8r.io/index"

var (
	// ErrNotExist mean this stack doesn't exist.
	ErrNotExist = errors.New("target stack doesn't exist")
)

// List all stacks
func List() ([]Stack, error) {
	type Index struct {
		Stacks []Stack `yaml:"stacks"`
	}
	b, err := getIndexYaml()
	if err != nil {
		return nil, err
	}
	idx := &Index{}
	if err := yaml.Unmarshal(b, idx); err != nil {
		return nil, err
	}
	return idx.Stacks, nil
}

func getIndexYaml() ([]byte, error) {
	var client http.Client
	resp, err := client.Get(StacksIndexURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("non-200 http code when fetching index contents")
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

// New returns a Stack object.
func New(name string) (*Stack, error) {
	const defaultVersion = "latest"

	url := fmt.Sprintf("https://stack.h8r.io/%s-%s.tar.gz", name, defaultVersion)
	s := &Stack{
		Path:    filepath.Join(state.GetCache(), name),
		Name:    name,
		URL:     url,
		Version: "0.0.1",
	}

	return s, nil
}

// Update upgrades the stack if necessary.
func (s *Stack) Update() error {
	ok := s.check()
	if !ok {
		s.clean()
		if err := s.pull(); err != nil {
			s.clean()
			return err
		}
	}
	return nil
}

// check checks if the stack is up to date.
func (s *Stack) check() bool {
	return false
}

func (s *Stack) pull() error {
	req := &getter.Request{
		Src: s.URL,
		Dst: filepath.Dir(s.Path),
	}
	err := util.GetWithTracker(req)
	if err != nil {
		return fmt.Errorf("failed to pull stack, please check stack name: %w", err)
	}
	return nil
}

func (s *Stack) clean() {
	if err := os.RemoveAll(s.Path); err != nil {
		panic(err)
	}
}
