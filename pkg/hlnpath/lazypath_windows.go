//go:build windows

package hlnpath

import "os"

func dataHome(cmd string) string {
	return filepath.Join(os.Getenv("APPDATA"), cmd, "data")
}

func configHome(cmd string) string {
	return filepath.Join(os.Getenv("APPDATA"), cmd, "config")
}

func cacheHome() string {
	return filepath.Join(os.Getenv("APPDATA"), cmd, "cache")
}
