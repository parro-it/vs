package writefs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
)

// WriteFS ...
type WriteFS interface {
	fs.FS
	OpenFile(name string, flag int, perm fs.FileMode) (FileWriter, error)
}

// FileWriter ...
type FileWriter interface {
	fs.File
	io.Writer
}

// ReadOnlyWriteFile ...
type ReadOnlyWriteFile struct {
	fs.File
}

func (f ReadOnlyWriteFile) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("file does not support write: %w", fs.ErrInvalid)
}

// OpenFile ...
func OpenFile(fsInst fs.FS, name string, flag int, perm fs.FileMode) (FileWriter, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{}
	}

	if fs, ok := fsInst.(WriteFS); ok {
		return fs.OpenFile(name, flag, perm)
	}

	if flag == os.O_RDONLY {
		file, err := fsInst.Open(name)
		if err != nil {
			return nil, err
		}
		return ReadOnlyWriteFile{file}, nil
	}

	return nil, fmt.Errorf("file system does not support write: %w", fs.ErrInvalid)
}

// WriteFile ...
func WriteFile(fsys fs.FS, name string, buf []byte) (n int, err error) {
	file, err := OpenFile(fsys, name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.FileMode(0644))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err = file.Write(buf)
	return
}
