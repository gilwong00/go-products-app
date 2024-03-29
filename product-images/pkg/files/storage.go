package files

import "io"

// Storage defines the behavior for file operations
type Storage interface {
	Save(path string, file io.Reader) error
}
