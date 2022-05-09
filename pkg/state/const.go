package state

import "path/filepath"

const (
	// HeighlinerNs Heighliner namespace
	HeighlinerNs = "heighliner"

	tfProviderConfigMapKey = "tf-provider"
	stackOutput            = "output.yaml"
	configTypeKey          = "heighliner.dev/config-type"
)

var (
	appInfo      = filepath.Join(".hln", "output.yaml")
	providerInfo = filepath.Join(".hln", "provider.tf")
)
