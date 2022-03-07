package state

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path"
	"sync"

	"github.com/hashicorp/go-getter/v2"
	"github.com/otiai10/copy"
	"github.com/rs/zerolog/log"
)

// Stack is a CloudNative app template
type Stack struct {
	Name        string `json:"name" yaml:"name"`               // The Name of the stack
	Path        string `json:"path"`                           // Path to localstorage
	URL         string `json:"url"`                            // URL to pull the stack
	Version     string `json:"version" yaml:"version"`         // Version
	Description string `json:"description" yaml:"description"` // A short description
}

var (
	// ErrStackNotExist means heighliner can't find the stack in localstorage
	ErrStackNotExist = errors.New("target stack doesn't exist")
)

var (
	// HeighlinerCacheHome is the dir where stacks are stored locally
	HeighlinerCacheHome string
)

var (
	// SampleStack is a demo stack that echoes your input
	SampleStack = &Stack{
		Name:        "sample",
		URL:         "https://stack.h8r.io/sample-latest.tar.gz",
		Description: "Sample is a light-weight stack mainly used for test",
		Version:     "1.0.0",
	}
	// GoGinStack is a go microservice architecture app
	GoGinStack = &Stack{
		Name:        "go-gin-stack",
		URL:         "https://stack.h8r.io/go-gin-stack-latest.tar.gz",
		Description: "Go-gin-stack helps you configure many cloud native components including prometheus, grafana, nocalhost, etc.",
		Version:     "1.0.0",
	}
)

// Stacks stores all stacks that currently usable
var Stacks = map[string]*Stack{
	"sample":       SampleStack,
	"go-gin-stack": GoGinStack,
}

// NewStack returns a Stack struct
func NewStack(name string) (*Stack, error) {
	HeighlinerCacheHome = os.Getenv("HEIGHLINER_CACHE_HOME")
	if HeighlinerCacheHome == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user cache dir: %w", err)
		}
		HeighlinerCacheHome = path.Join(cacheDir, "heighliner")
	}
	return &Stack{
		Name: name,
		Path: path.Join(HeighlinerCacheHome, "repository"),
	}, nil
}

// Pull downloads and decompresses a stack
func (s *Stack) Pull(url string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working dir: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := &getter.Client{}
	req := &getter.Request{
		Src:              url,
		Dst:              s.Path,
		Pwd:              pwd,
		ProgressListener: defaultProgressBar,
	}

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
		wg.Wait()
		log.Info().Msgf("signal %s", sig)
		return nil
	case <-ctx.Done():
		wg.Wait()
		log.Info().Msgf("successfully pull stack %s", s.Name)
		return nil
	case err := <-errChan:
		wg.Wait()
		return fmt.Errorf("failed to pull stack %s: %w", s.Name, err)
	}
}

// Check checks the status of target stack
func (s *Stack) Check() error {
	dir := path.Join(s.Path, s.Name)

	_, err := os.Stat(dir)
	if err != nil {
		return ErrStackNotExist
	}

	return nil
}

// Copy copies the stack into dest
func (s *Stack) Copy(dest string) error {
	src := path.Join(s.Path, s.Name)
	err := copy.Copy(src, dest)
	if err != nil {
		return fmt.Errorf("failed to copy stack %s: %w", s.Name, err)
	}
	return nil
}
