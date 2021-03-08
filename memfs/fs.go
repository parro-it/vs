package memfs

import (
	"io/fs"
	"os"
	"path"
	"strings"
	"testing/fstest"
	"time"

	"github.com/parro-it/vs/writefs"
)

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
			if os.IsNotExist(err) {
				_, err = fsys.OpenFile(curr, flag, perm)
				if err != nil {
					return nil, err
				}

			}
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
		delete(fsys.MapFS, name)
		return nil, nil
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
		fsys.MapFS[name] = file

	}

	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}

	return memWriteFile{
		File: f,
		file: file,
	}, nil
}
