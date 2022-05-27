package schema

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/mitchellh/go-homedir"
	"sigs.k8s.io/yaml"
)

var (
	// ErrNotExist means no input schema for interactive prompt.
	ErrNotExist = errors.New("no schema found in current satck")
)

// Schema represents a input schema of a stack.
type Schema struct {
	// Dir is the path to stack! Not schema directly.
	Dir        string
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
func New(dir string) *Schema {
	return &Schema{
		Dir: dir,
	}
}

// Show displays the info of input shcema.
func (s Schema) Show(o io.Writer) {

	w := tabwriter.NewWriter(o, 0, 4, 2, ' ', tabwriter.TabIndent)
	defer func() {
		if err := w.Flush(); err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Fprintf(w, "\nPARAMETERS LIST:\n")
	fmt.Fprintln(w, "PARAMETER\tTYPE\tKEY\tDEFAULT\tREQUIRED\tDESCRIPTION")
	for _, p := range s.Parameters {
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%v\t%s", p.Title, p.Type, p.Key, p.Default, p.Required, p.Description)
		fmt.Fprintln(w, line)
	}
}

// AutomaticEnv sets envs automatically.
func (s *Schema) AutomaticEnv(interactive bool) error {
	var err = s.LoadSchema()
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
				key, val := v.Key, v.Default
				val, err := homedir.Expand(val)
				if err != nil {
					return err
				}
				val, err = filepath.Abs(val)
				if err != nil {
					return err
				}
				if err := os.Setenv(key, val); err != nil {
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

// LoadSchema reads the schema
func (s *Schema) LoadSchema() error {
	b, err := ioutil.ReadFile(filepath.Join(s.Dir, "schemas", "schema.yaml"))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrNotExist, err.Error())
	}
	if err = yaml.Unmarshal(b, s); err != nil {
		return fmt.Errorf("syntax error in schema: %w", err)
	}
	return nil
}
