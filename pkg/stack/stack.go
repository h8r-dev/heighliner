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
	"sigs.k8s.io/yaml"

	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
)

const (
	// HlnRepoURL is official repo.
	HlnRepoURL = "https://stack.h8r.io"
	// MetaFileName is the name of metadata file.
	MetaFileName = "metadata.yaml"
)

// Stack is a CloudNative application template.
type Stack struct {
	Path string

	Name        string `json:"name" yaml:"name"`
	URL         string `json:"url" yaml:"url"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
}

// Metadata of the schema.
type Metadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	OwnerRef    Owner  `json:"owner"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	URL         string `json:"url"`
	Tags        []*Tag `json:"tags"`
}

// Owner info.
type Owner struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

// Tag is a key word.
type Tag string

// LoadMeta loads metadata file.
func LoadMeta(path string) (*Metadata, error) {
	metafile := filepath.Join(path, MetaFileName)
	b, err := os.ReadFile(metafile)
	if err != nil {
		return nil, err
	}
	stackMeta := &Metadata{}
	if err := yaml.Unmarshal(b, stackMeta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mettadata: %w", err)
	}
	return stackMeta, nil
}

// Show displays the info of metadata.
func (m Metadata) Show(w io.Writer) {
	fmt.Fprintf(w, "\nNAME: %s\n", m.Name)
	fmt.Fprintf(w, "VERSION: %s\n", m.Version)
	fmt.Fprintf(w, "DESCRIPTION: %s\n", m.Description)
}

// List all stacks
func List() ([]Stack, error) {
	b, err := getIndexYaml(HlnRepoURL)
	if err != nil {
		return nil, err
	}
	idx := &struct {
		Stacks []Stack `yaml:"stacks"`
	}{}
	if err := yaml.Unmarshal(b, idx); err != nil {
		return nil, err
	}
	return idx.Stacks, nil
}

func getIndexYaml(repoURL string) ([]byte, error) {
	var client http.Client
	indexFile := "index.yaml"
	indexURL := repoURL + "/" + indexFile
	resp, err := client.Get(indexURL)
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
func New(name, version string) (*Stack, error) {
	const defaultVersion = "latest"
	if version == "" {
		version = defaultVersion
	}
	url := fmt.Sprintf("https://stack.h8r.io/%s-%s.tar.gz", name, version)
	s := &Stack{
		Path:    filepath.Join(state.GetCache(), name),
		Name:    name,
		URL:     url,
		Version: version,
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
