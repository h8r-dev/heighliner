package state

import (
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

var (
	// HeighlinerCacheHome is the dir where stacks are stored locally
	HeighlinerCacheHome string
)

func init() {
	if err := initHeighlinerCache(); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize localstorage")
	}
}

func initHeighlinerCache() error {
	HeighlinerCacheHome = os.Getenv("HEIGHLINER_CACHE_HOME")
	if HeighlinerCacheHome == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return fmt.Errorf("failed to get user cache dir: %w", err)
		}
		HeighlinerCacheHome = path.Join(cacheDir, "heighliner")
	}
	err := os.MkdirAll(HeighlinerCacheHome, 0755)
	if err != nil {
		return fmt.Errorf("failed to create dir %s: %w", HeighlinerCacheHome, err)
	}
	return nil
}

// CleanHeighlinerCaches cleans all cached cuemods and stacks
func CleanHeighlinerCaches() error {
	err := os.RemoveAll(HeighlinerCacheHome)
	if err != nil {
		return err
	}
	return nil
}
