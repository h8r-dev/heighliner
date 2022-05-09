package state

import (
	"fmt"
	"github.com/h8r-dev/heighliner/pkg/state/app"
	"github.com/otiai10/copy"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type LocalFileState struct {
}

func (l *LocalFileState) LoadOutput(appName string) (*app.Output, error) {
	b, err := os.ReadFile(filepath.Join(".hln", "output.yaml"))
	if err != nil {
		return nil, err
	}
	output := &app.Output{}
	err = yaml.Unmarshal(b, output)
	return output, err
}

func (l *LocalFileState) LoadTFProvider(appName string) (string, error) {
	return "", nil
}

func (l *LocalFileState) SaveOutputAndTFProvider(appName string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := copy.Copy(stackOutput, filepath.Join(pwd, appInfo)); err != nil {
		return err
	}
	if err := os.Remove(stackOutput); err != nil {
		return err
	}
	ao, err := app.Load(filepath.Join(pwd, appInfo))
	if err != nil {
		return fmt.Errorf("failed to load app output: %w", err)
	}
	return copy.Copy(ao.SCM.TfProvider, filepath.Join(pwd, providerInfo))
}
