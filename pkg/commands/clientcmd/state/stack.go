package state

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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
	src := filepath.Join(s.Path, s.Name+".tar.gz")

	err := decompress(src, s.Path)
	if err != nil {
		return err
	}

	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}

type dirInfo struct {
	Name    string
	ModTime time.Time
}

func makeTarReader(filename string) (*tar.Reader, func(), error) {
	srcFile, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	content, err := ioutil.ReadAll(srcFile)
	if err != nil {
		srcFile.Close()
		return nil, nil, err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err = w.Write(content)
	if err != nil {
		srcFile.Close()
		return nil, nil, err
	}
	w.Close()
	gr, err := gzip.NewReader(&b)
	if err != nil {
		srcFile.Close()
		return nil, nil, err
	}

	closeFunc := func() {
		srcFile.Close()
		gr.Close()
	}
	tr := tar.NewReader(gr)
	return tr, closeFunc, nil
}

// decompress decompresses a tar.gz file into dest dir.
func decompress(tarFile, dest string) error {
	tr, closeFDs, err := makeTarReader(tarFile)
	if err != nil {
		return err
	}
	defer closeFDs()
	if dest != "" {
		_, err = makeDir(dest)
		if err != nil {
			return err
		}
	}
	currentDir := dirInfo{}

	// iterate until all files are decompresses
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				if currentDir.Name != "" {
					remodifyTime(currentDir.Name, currentDir.ModTime)
				}
				break
			} else {
				return err
			}
		}
		fi := header.FileInfo()
		fileName := filepath.Join(dest, header.Name)
		if !strings.HasPrefix(fileName, currentDir.Name) {
			remodifyTime(currentDir.Name, currentDir.ModTime)
		}
		if fi.IsDir() {
			foldName, err := makeDir(fileName)
			if err != nil {
				return err
			}
			currentDir = dirInfo{
				foldName,
				fi.ModTime(),
			}
			continue
		}
		file, err := createFile(fileName)
		if err != nil {
			return fmt.Errorf("can not create file %v: %v", fileName, err)
		}
		_, err = io.Copy(file, tr)
		if err != nil {
			return err
		}
		file.Close()
		remodifyTime(fileName, header.ModTime)
	}
	return nil
}

func remodifyTime(name string, modTime time.Time) {
	if name == "" {
		return
	}
	atime := time.Now()
	_ = os.Chtimes(name, atime, modTime)
}

func makeDir(name string) (string, error) {
	if name != "" {
		_, err := os.Stat(name)
		if err != nil {
			err = os.MkdirAll(name, 0755)
			if err != nil {
				return "", fmt.Errorf("can not make directory: %v", err)
			}
			return name, nil
		}
		return "", nil
	}
	return "", fmt.Errorf("can not make directory without a name: %v", name)
}

func createFile(name string) (*os.File, error) {
	dir := path.Dir(name)
	if dir != "" {
		_, err := os.Lstat(dir)
		if err != nil {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return nil, err
			}
		}
	}
	return os.Create(name)
}
