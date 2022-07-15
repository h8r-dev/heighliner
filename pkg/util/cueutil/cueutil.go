package cueutil

import (
	"fmt"
	"os"

	"cuelang.org/go/cue/format"
	"cuelang.org/go/encoding/yaml"
)

// ConvertYamlToCue and write cue file to dest
func ConvertYamlToCue(from string, to string) error {
	a, err := yaml.Extract(from, nil)
	if err != nil {
		return fmt.Errorf("failed to parse yaml file: %w", err)
	}
	b, err := format.Node(a, format.Simplify())
	if err != nil {
		return fmt.Errorf("failed to convert yaml into cue: %w", err)
	}
	if _, err := os.Create(to); err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	b = append([]byte("package plans\n\n"), b...)
	if err := os.WriteFile(to, b, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
