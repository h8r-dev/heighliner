package schema

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	// ErrNoSchema means no input schema for interactive prompt.
	ErrNoSchema = errors.New("no schema found in current satck")
)

// Schema represents a input schema of a stack.
type Schema struct {
	Parameters []Parameter `yaml:"parameters"`
}

// Parameter is a field in the schema.
type Parameter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Key         string `yaml:"key"`
	Value       string `yaml:"value"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`
}

// New creates and returns a schema.
func New() *Schema {
	return &Schema{}
}

// AutomaticEnv sets envs automatically.
func (s *Schema) AutomaticEnv(interactive bool) error {
	var err = s.load()
	if err != nil {
		return err
	}

	for _, v := range s.Parameters {
		// Try to fetch value from env
		val := os.Getenv(v.Key)
		if val != "" {
			continue
		}
		// Promt interactively or look for default values
		if interactive {
			if err := startUI(v); err != nil {
				return err
			}
		} else {
			switch {
			case v.Default != "":
				if err := os.Setenv(v.Key, v.Default); err != nil {
					panic(err)
				}
				continue
			case !v.Required:
				continue
			default:
				return fmt.Errorf("couldn't find value of %s, which is required", v.Title)
			}
		}
	}

	return nil
}

// NOTE: Make sure you are already in the project dir.
func (s *Schema) load() error {
	file, err := os.Open(path.Join("schemas", "schema.yaml"))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrNoSchema, err.Error())
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read schema file %w", err)
	}
	if err = yaml.Unmarshal(b, s); err != nil {
		return fmt.Errorf("syntax error in schema: %w", err)
	}
	return nil
}
