package state

import (
	"github.com/h8r-dev/heighliner/pkg/state/app"
)

// State Heighliner application state
type State interface {
	ListApps() ([]string, error)
	LoadOutput(appName string) (*app.Output, error)
	LoadTFProvider(appName string) (string, error)
	SaveOutputAndTFProvider(appName string) error
}
