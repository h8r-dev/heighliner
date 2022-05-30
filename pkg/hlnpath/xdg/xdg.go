// Package xdg holds constants pertaining to XDG Base Directory Specification.
//
// The XDG Base Directory Specification https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
// specifies the environment variables that define user-specific base directories for various categories of files.
package xdg

const (
	// CacheHomeEnvVar is the environment variable used by the
	// XDG base directory specification for the cache directory.
	CacheHomeEnvVar = "XDG_CACHE_HOME"

	// ConfigHomeEnvVar is the environment variable used by the
	// XDG base directory specification for the config directory.
	ConfigHomeEnvVar = "XDG_CONFIG_HOME"

	// DataHomeEnvVar is the environment variable used by the
	// XDG base directory specification for the data directory.
	DataHomeEnvVar = "XDG_DATA_HOME"
)
