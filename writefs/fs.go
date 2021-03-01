package writefs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"testing/fstest"
	"time"
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

type readOnlyWriteFile struct {
	fs.File
}

func (f readOnlyWriteFile) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("file does not support write: %w", fs.ErrInvalid)
}

// OpenFile ...
func OpenFile(fsInst fs.FS, name string, flag int, perm fs.FileMode) (FileWriter, error) {
	if fs, ok := fsInst.(WriteFS); ok {
		return fs.OpenFile(name, flag, perm)
	}

	if flag == os.O_RDONLY {
		file, err := fsInst.Open(name)
		if err != nil {
			return nil, err
		}
		return readOnlyWriteFile{
			File: file,
		}, nil
	}

	return nil, fmt.Errorf("file system does not support write: %w", fs.ErrInvalid)
}

type osWriteFS struct {
	fs.FS
	root string
}

// OpenFile ...
func (fsinst osWriteFS) OpenFile(name string, flag int, perm fs.FileMode) (FileWriter, error) {
	realPath := path.Join(fsinst.root, name)
	return os.OpenFile(realPath, flag, perm)
}

// MapWriteFS ...
type MapWriteFS struct {
	fstest.MapFS
}

type memWriteFile struct {
	fs.File
	file *fstest.MapFile
}

func (f memWriteFile) Write(buf []byte) (n int, err error) {
	f.file.Data = append(f.file.Data, buf...)
	return len(buf), nil
}

// OpenFile ...
func (fsinst MapWriteFS) OpenFile(name string, flag int, perm fs.FileMode) (FileWriter, error) {

	file, exists := fsinst.MapFS[name]

	if flag == os.O_RDONLY {
		if !exists {
			return nil, fs.ErrNotExist
		}
		f, err := fsinst.Open(name)
		if err != nil {
			return nil, err
		}
		return readOnlyWriteFile{f}, nil
	}

	if exists {
		if flag&os.O_TRUNC == os.O_TRUNC {
			file.Data = []byte{}
		} else if flag&os.O_EXCL == os.O_EXCL {
			return nil, fs.ErrExist
		} else if flag&os.O_APPEND == 0 {
			// non append open of existing files is not supported
			return nil, fs.ErrInvalid
		}
	} else {
		if flag&os.O_CREATE == 0 {
			return nil, fs.ErrNotExist
		}

		file = &fstest.MapFile{
			Data:    []byte{},
			Mode:    perm,
			ModTime: time.Now(),
		}
		fsinst.MapFS[name] = file

	}

	f, err := fsinst.Open(name)
	if err != nil {
		return nil, err
	}

	return memWriteFile{
		File: f,
		file: file,
	}, nil
}

// DirWriteFS ...
func DirWriteFS(dir string) WriteFS {
	return osWriteFS{
		FS:   os.DirFS(dir),
		root: dir,
	}
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
