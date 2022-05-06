package cmd

import (
	"path/filepath"
)

const (
	stackOutput            = "output.yaml"
	heighlinerNs           = "heighliner"
	buildKitName           = "buildkitd"
	terraformDir           = ".hln"
	tfProviderConfigMapKey = "tf-provider"
	configTypeKey          = "heighliner.dev/config-type"
)

var (
	appInfo      = filepath.Join(".hln", "output.yaml")
	providerInfo = filepath.Join(".hln", "provider.tf")
)
