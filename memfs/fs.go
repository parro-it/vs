package memfs

import (
	"io/fs"
	"os"
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
func (fsinst MapWriteFS) OpenFile(name string, flag int, perm fs.FileMode) (writefs.FileWriter, error) {

	file, exists := fsinst.MapFS[name]

	if flag == os.O_RDONLY {
		if !exists {
			return nil, fs.ErrNotExist
		}
		f, err := fsinst.Open(name)
		if err != nil {
			return nil, err
		}
		return writefs.ReadOnlyWriteFile{File: f}, nil
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
