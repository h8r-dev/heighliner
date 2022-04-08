package state

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

// GetCache returns the path to the cache directory.
func GetCache() string {
	var cache string
	cacheDir := viper.GetString("cache_home")
	if cacheDir == "" {
		userCache, err := os.UserCacheDir()
		if err != nil {
			panic(err)
		}
		cache = path.Join(userCache, "heighliner")
	} else {
		cache = path.Join(cacheDir, "heighliner")
	}
	return cache
}

// CleanCache removes all cached resources.
func CleanCache() {
	if err := os.RemoveAll(GetCache()); err != nil {
		panic(err)
	}
}
