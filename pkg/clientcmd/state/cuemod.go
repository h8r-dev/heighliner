package state

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/hashicorp/go-getter/v2"
	"github.com/otiai10/copy"
)

// CueMod represents a Cue module
type CueMod struct {
	Name   string `json:"name" yaml:"name"`
	Path   string `json:"path"`
	URL    string `json:"url"`
	Subdir string `json:"subdir"`
}

var (
	// ErrCueModNotExist means heighliner can't find the cuemod in localstorage
	ErrCueModNotExist = errors.New("target cuemod doesn't exist")
)

var (
	heighlinerCueLib = &CueMod{
		Name: "cuelib",
		URL:  "https://github.com/h8r-dev/cuelib.git",
	}
)

// NewCueMod creates a *CueMod struct and returns it
func NewCueMod(name, src string) (*CueMod, error) {
	HeighlinerCacheHome = os.Getenv("HEIGHLINER_CACHE_HOME")
	if HeighlinerCacheHome == "" {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user cache dir: %w", err)
		}
		HeighlinerCacheHome = path.Join(cacheDir, "heighliner")
	}
	pr, err := url.Parse(src)
	if err != nil {
		return nil, err
	}
	return &CueMod{
		Name:   name,
		URL:    src,
		Path:   path.Join(HeighlinerCacheHome, "cuemod"),
		Subdir: path.Join(pr.Host, path.Dir(pr.Path)),
	}, nil
}

// PrePareCueMod pull the cue mod from github repo if necessary and copy it into dst dir
func PrePareCueMod(dst string) error {
	c, err := NewCueMod(heighlinerCueLib.Name, heighlinerCueLib.URL)
	if err != nil {
		return err
	}

	err = c.Check()
	if err != nil && errors.Is(err, ErrCueModNotExist) {
		err := c.Pull()
		if err != nil {
			return fmt.Errorf("failed to pull cuemod:%s,%w", c.Name, err)
		}
	} else if err != nil {
		return err
	}

	err = c.Copy(dst)
	if err != nil {
		return err
	}

	return nil
}

// Pull downloads a cuemod from it's github repo
func (c *CueMod) Pull() error {
	req := &getter.Request{
		Src: "git::" + c.URL,
		Dst: path.Join(c.Path, c.Subdir, c.Name),
	}
	err := getWithTracker(req)
	if err != nil {
		return fmt.Errorf("failed to download stack: %w", err)
	}
	return nil
}

// Check checks if a cuemod exists in localstorage
func (c *CueMod) Check() error {
	dir := path.Join(c.Path, c.Subdir, c.Name)

	_, err := os.Stat(dir)
	if err != nil {
		return ErrCueModNotExist
	}

	return nil
}

// Copy copies the CueMod into dst
func (c *CueMod) Copy(dst string) error {
	realDst := path.Join(dst, "cue.mod", "pkg")
	err := copy.Copy(c.Path, realDst)
	if err != nil {
		return fmt.Errorf("failed to copy cuemod %s: %w", c.Name, err)
	}
	return nil
}
