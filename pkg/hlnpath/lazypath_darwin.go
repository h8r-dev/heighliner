//go:build darwin

package hlnpath

import (
	"path/filepath"

	"k8s.io/client-go/util/homedir"
)

func dataHome() string {
	return filepath.Join(homedir.HomeDir(), "Library")
}

func configHome() string {
	return filepath.Join(homedir.HomeDir(), "Library", "Preferences")
}

func cacheHome() string {
	return filepath.Join(homedir.HomeDir(), "Library", "Caches")
}
