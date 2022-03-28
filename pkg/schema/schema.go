package schema

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

// Schema represents a input schema of a stack
type Schema struct {
	Parameters []Parameter `yaml:"parameters"`
}

// Parameter is a field in the schema
type Parameter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Type        string `yaml:"type"`
	Key         string `yaml:"key"`
	Value       string `yaml:"value"`
	Default     string `yaml:"default"`
}

// New creates and returns a schema
func New() *Schema {
	return &Schema{}
}

// Load loads the schema from src file
func (s *Schema) Load() error {
	file, err := os.Open(path.Join("schemas", "schema.yaml"))
	if err != nil {
		return fmt.Errorf("couldn't find schema :%w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(b, s); err != nil {
		return fmt.Errorf("syntax error in schema: %w", err)
	}
	return nil
}

// SetEnv sets up input values
func (s *Schema) SetEnv(m map[string]interface{}) error {
	for _, v := range s.Parameters {
		val, ok := m[v.Key]
		if ok && val != nil {
			if err := os.Setenv(v.Key, val.(string)); err != nil {
				panic(err)
			}
			continue
		}
		val = os.Getenv(v.Key)
		if val == "" && v.Default != "" {
			if err := os.Setenv(v.Key, v.Default); err != nil {
				panic(err)
			}
			continue
		}
		return fmt.Errorf("couldn't find value of %s, which is required", v.Title)
	}
	return nil
}
