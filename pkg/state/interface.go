package state

import "github.com/h8r-dev/heighliner/pkg/state/app"

type State interface {
	LoadOutput(appName string) (*app.Output, error)
	LoadTfProvider(appName string) (string, error)
	SaveOutputAndTfProvider(appName string) error
}
