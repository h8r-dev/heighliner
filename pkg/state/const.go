package state

import "path/filepath"

const (
	// InfraNs is the namespace of infra
	InfraNs = "heighliner-infra"
	// InfraConfigMap Infra configmap
	InfraConfigMap = "heighliner-infra-config"
	// InfraEntry Infra entry
	InfraEntry = "infra"
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
