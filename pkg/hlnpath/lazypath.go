package hlnpath

import (
	"os"
	"path/filepath"

	"github.com/h8r-dev/heighliner/pkg/hlnpath/xdg"
)

const (
	// CacheHomeEnvVar is the environment variable used by Hln
	// for the cache directory. When no value is set a default is used.
	CacheHomeEnvVar = "HLN_CACHE_HOME"

	// ConfigHomeEnvVar is the environment variable used by Hln
	// for the config directory. When no value is set a default is used.
	ConfigHomeEnvVar = "HLN_CONFIG_HOME"

	// DataHomeEnvVar is the environment variable used by Hln
	// for the data directory. When no value is set a default is used.
	DataHomeEnvVar = "HLN_DATA_HOME"
)

// lazypath is an lazy-loaded path buffer for the XDG base directory specification.
type lazypath string

func (l lazypath) path(hlnEnvVar, xdgEnvVar string, defaultFn func(cmd string) string, elem ...string) string {

	// There is an order to checking for a path.
	// 1. See if a Hln specific environment variable has been set.
	// 2. Check if an XDG environment variable is set
	// 3. Fall back to a default
	base := os.Getenv(hlnEnvVar)
	if base != "" {
		return filepath.Join(base, filepath.Join(elem...))
	}
	base = os.Getenv(xdgEnvVar)
	if base != "" {
		return filepath.Join(base, string(l), filepath.Join(elem...))
	}
	base = defaultFn(string(l))
	return filepath.Join(base, filepath.Join(elem...))
}

// cachePath defines the base directory relative to which user specific non-essential data files
// should be stored.
func (l lazypath) cachePath(elem ...string) string {
	return l.path(CacheHomeEnvVar, xdg.CacheHomeEnvVar, cacheHome, filepath.Join(elem...))
}

// configPath defines the base directory relative to which user specific configuration files should
// be stored.
func (l lazypath) configPath(elem ...string) string {
	return l.path(ConfigHomeEnvVar, xdg.ConfigHomeEnvVar, configHome, filepath.Join(elem...))
}

// dataPath defines the base directory relative to which user specific data files should be stored.
func (l lazypath) dataPath(elem ...string) string {
	return l.path(DataHomeEnvVar, xdg.DataHomeEnvVar, dataHome, filepath.Join(elem...))
}
