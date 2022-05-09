package cmd

import (
	"path/filepath"
)

const (
	buildKitName = "buildkitd"
	terraformDir = ".hln"
)

var (
	appInfo      = filepath.Join(".hln", "output.yaml")
	providerInfo = filepath.Join(".hln", "provider.tf")
)
