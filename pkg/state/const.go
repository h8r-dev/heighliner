package state

import "path/filepath"

const (
	// HeighlinerNs Heighliner namespace
	InfraNs        = "heighliner-infra"
	InfraConfigMap = "heighliner-infra-config"
	InfraEntry     = "infra"
	HeighlinerNs   = "heighliner"

	tfProviderConfigMapKey = "tf-provider"
	stackOutput            = "output.yaml"
	configTypeKey          = "heighliner.dev/config-type"
)

var (
	appInfo      = filepath.Join(".hln", "output.yaml")
	providerInfo = filepath.Join(".hln", "provider.tf")
)
