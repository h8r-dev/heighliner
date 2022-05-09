package state

import "path/filepath"

const (
	tfProviderConfigMapKey = "tf-provider"
	HeighlinerNs           = "heighliner"
	stackOutput            = "output.yaml"
	configTypeKey          = "heighliner.dev/config-type"
)

var (
	appInfo      = filepath.Join(".hln", "output.yaml")
	providerInfo = filepath.Join(".hln", "provider.tf")
)
