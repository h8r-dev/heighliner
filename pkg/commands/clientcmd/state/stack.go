package state

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type Stack struct {
	Name        string         `json:"name"`
	Path        string         `json:"path"`
	Url         string         `json:"url"`
	Version     string         `json:"version"`
	Description string         `json:"description"`
	Inputs      []*InputSchema `json:"inputSchema"`
}

func (s *Stack) Download() error {
	fp := filepath.Join(s.Path, s.Name+".tar.gz")
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	rsp, err := http.Get(s.Url)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *Stack) Decompress() error {
	command := exec.Command("tar",
		"-zxvf", filepath.Join(s.Path, s.Name+".tar.gz"),
		"-C", s.Path)
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}
	err := command.Run()
	if err != nil {
		return err
	}
	fmt.Println(command.Stdout.(*bytes.Buffer).String())

	err = os.Remove(filepath.Join(s.Path, s.Name+".tar.gz"))
	if err != nil {
		return err
	}

	return nil
}
