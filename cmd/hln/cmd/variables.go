package cmd

import (
	"path/filepath"
)

const (
	stackOutput  = "output.yaml"
	heighlinerNs = "heighliner"
	buildKitName = "buildkitd"
)

var (
	appInfo = filepath.Join(".hln", "output.yaml")
)
