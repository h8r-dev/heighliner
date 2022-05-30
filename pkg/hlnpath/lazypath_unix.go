//go:build !windows && !darwin

package hlnpath

import (
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

// dataHome defines the base directory relative to which user specific data files should be stored.
//
// If $XDG_DATA_HOME is either not set or empty, a default equal to $HOME/.local/share is used.
func dataHome() string {
	return filepath.Join(homedir.HomeDir(), ".local", "share")
}

// configHome defines the base directory relative to which user specific configuration files should
// be stored.
//
// If $XDG_CONFIG_HOME is either not set or empty, a default equal to $HOME/.config is used.
func configHome() string {
	return filepath.Join(homedir.HomeDir(), ".config")
}

// cacheHome defines the base directory relative to which user specific non-essential data files
// should be stored.
//
// If $XDG_CACHE_HOME is either not set or empty, a default equal to $HOME/.cache is used.
func cacheHome() string {
	return filepath.Join(homedir.HomeDir(), ".cache")
}
