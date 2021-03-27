package memfs

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing/fstest"
	"time"

	"github.com/parro-it/vs/writefs"
)

// MapWriteFS ...
type MapWriteFS struct {
	fstest.MapFS
}

// New ...
func New() *MapWriteFS {
	return &MapWriteFS{fstest.MapFS{}}
}

// NewFS ...
func NewFS() writefs.WriteFS {
	return &MapWriteFS{fstest.MapFS{}}
}

type memWriteFile struct {
	fs.File
	file   *fstest.MapFile
	cursor int
}

func (f *memWriteFile) Write(buf []byte) (n int, err error) {
	sz := len(buf)
	if f.cursor == len(f.file.Data) {
		f.file.Data = append(f.file.Data, buf...)
	} else {
		buf = append(buf, f.file.Data[f.cursor+len(buf):]...)
		f.file.Data = append(f.file.Data[:f.cursor], buf...)
	}
	f.cursor += sz
	return sz, nil
}

// OpenFile ...
func (fsys MapWriteFS) OpenFile(name string, flag int, perm fs.FileMode) (writefs.FileWriter, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{}
	}

	file, exists := fsys.MapFS[name]

	if flag == os.O_RDONLY {
		if !exists {
			return nil, fs.ErrNotExist
		}
		f, err := fsys.Open(name)
		if err != nil {
			return nil, err
		}
		return writefs.ReadOnlyWriteFile{File: f}, nil
	}

	if flag&os.O_CREATE == os.O_CREATE && perm.IsDir() {
		if exists {
			return nil, fs.ErrExist
		}
		segments := strings.Split(name, "/")
		curr := "."
		for _, seg := range segments[:len(segments)-1] {
			curr = path.Join(curr, seg)
			_, err := fsys.Stat(curr)
			/*if os.IsNotExist(err) {
				_, err = fsys.OpenFile(curr, flag, perm)
				if err != nil {
					return nil, err
				}

			}*/
			if err != nil {
				return nil, err
			}
		}
		fsys.MapFS[name] = &fstest.MapFile{
			Mode:    perm,
			ModTime: time.Now(),
		}
		return nil, nil
	}

	if flag == os.O_TRUNC {
		if !exists {
			return nil, fs.ErrNotExist
		}
		if file.Mode.IsDir() {
			files, err := fs.ReadDir(fsys, name)
			if err != nil {
				return nil, err
			}
			if len(files) != 0 {
				return nil, fmt.Errorf("%w: directory `%s` not empty", fs.ErrInvalid, name)
			}
		}
		delete(fsys.MapFS, name)
		return nil, nil
	}
	cursor := 0
	if exists {
		if flag&os.O_TRUNC == os.O_TRUNC {
			file.Data = []byte{}
		} else if flag&os.O_EXCL == os.O_EXCL {
			return nil, fs.ErrExist
		} else if flag&os.O_APPEND == os.O_APPEND {
			cursor += len(file.Data)
		}
	} else {
		if flag&os.O_CREATE == 0 {
			return nil, fs.ErrNotExist
		}
		if name != "." {
			// check that parent exists
			parent, err := fs.Stat(fsys, filepath.Dir(name))
			if err != nil {
				return nil, err
			}
			if !parent.IsDir() {
				return nil, fmt.Errorf("parent directory `%s` is a file: %w", filepath.Dir(name), fs.ErrInvalid)
			}
		}

		file = &fstest.MapFile{
			Data:    []byte{},
			Mode:    perm,
			ModTime: time.Now(),
		}
		fsys.MapFS[name] = file

	}

	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}

	return &memWriteFile{
		File:   f,
		file:   file,
		cursor: cursor,
	}, nil
}
