package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Local is an implementation of the Storage interface which works with the
// local disk on the current machine
type Local struct {
	maxFileSize int // maximum number of bytes for files
	basePath    string
}

func NewLocal(basePath string, maxSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	return &Local{basePath: p}, nil
}

// path string, contents io.Reader
func (l *Local) Save(path string, contents io.Reader) error {
	p := l.fullPath(path)
	// get the directory and make sure it exists
	d := filepath.Dir(p)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}

	// if file already exist, delete it
	_, err = os.Stat(p)
	if err == nil {
		err = os.Remove(p)
		if err != nil {
			return fmt.Errorf("unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		// if this is anything other than a not exists error
		return fmt.Errorf("unable to get file info: %w", err)
	}

	// create a new file at given path
	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer f.Close()

	// write contents to new file
	// make sure we honor max bytes
	// buffer reading and loop over and check if it exceeds max num of bytes
	_, err = io.Copy(f, contents)
	if err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}
	return nil
}

// returns the absolute path
func (l *Local) fullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(l.basePath, path)
}
