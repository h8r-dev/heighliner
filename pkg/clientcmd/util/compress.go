package util

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type dirInfo struct {
	Name    string
	ModTime time.Time
}

// Decompress decompresses a tar.gz file into dest dir.
func Decompress(tarFile, dest string) error {
	tr, err := makeTarReader(tarFile)
	if err != nil {
		return err
	}
	if dest != "" {
		_, err = makeDir(dest)
		if err != nil {
			return err
		}
	}
	currentDir := dirInfo{}

	// iterate until all files are decompressed
	for {
		header, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if currentDir.Name != "" {
					remodifyTime(currentDir.Name, currentDir.ModTime)
				}
				break
			} else {
				return err
			}
		}
		fi := header.FileInfo()
		fileName := path.Join(dest, header.Name)
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
			return fmt.Errorf("can not create file %v: %w", fileName, err)
		}
		_, err = io.Copy(file, tr)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		remodifyTime(fileName, header.ModTime)
	}
	return nil
}

func makeTarReader(filename string) (*tar.Reader, error) {
	srcFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return nil, err
	}
	err = srcFile.Close()
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err = w.Write(content)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(&b)
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gr)
	return tr, nil
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
			err = os.MkdirAll(name, 0750)
			if err != nil {
				return "", fmt.Errorf("can not make directory: %w", err)
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
			err := os.MkdirAll(dir, 0750)
			if err != nil {
				return nil, err
			}
		}
	}
	return os.Create(name)
}
