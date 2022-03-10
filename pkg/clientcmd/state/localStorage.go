package state

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"sync"

	"github.com/hashicorp/go-getter/v2"
	"github.com/rs/zerolog/log"
)

var (
	// HeighlinerCacheHome is the dir where stacks are stored locally
	HeighlinerCacheHome string
)

func initHeighlinerCache() error {
	HeighlinerCacheHome = os.Getenv("HEIGHLINER_CACHE_HOME")
	if HeighlinerCacheHome == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return fmt.Errorf("failed to get user cache dir: %w", err)
		}
		HeighlinerCacheHome = path.Join(cacheDir, "heighliner")
	}
	return nil
}

// CleanHeighlinerCaches cleans all cached cuemods and stacks
func CleanHeighlinerCaches() error {
	err := CleanStacks()
	if err != nil {
		return err
	}
	err = CleanCueMods()
	if err != nil {
		return err
	}
	return nil
}

func getWithTracker(req *getter.Request) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working dir: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := &getter.Client{}

	req.Pwd = pwd
	req.ProgressListener = defaultProgressBar

	wg := sync.WaitGroup{}
	wg.Add(1)
	errChan := make(chan error, 2)
	go func() {
		defer wg.Done()
		defer cancel()
		if _, err := client.Get(ctx, req); err != nil {
			errChan <- err
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		signal.Reset(os.Interrupt)
		cancel()
		log.Info().Msgf("signal %s", sig)
		return nil
	case <-ctx.Done():
		wg.Wait()
		return nil
	case err := <-errChan:
		wg.Wait()
		return err
	}
}
