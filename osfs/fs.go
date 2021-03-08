package osfs

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"syscall"

	"github.com/parro-it/vs/writefs"
)

type osWriteFS struct {
	fs.FS
	root string
}

// OpenFile ...
func (fsinst osWriteFS) OpenFile(name string, flag int, perm fs.FileMode) (writefs.FileWriter, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{}
	}
	realPath := path.Join(fsinst.root, name)

	if flag&os.O_CREATE == os.O_CREATE && perm.IsDir() {
		err := os.Mkdir(realPath, perm)
		if err != nil {
			if os.IsExist(err) {
				return nil, err
			}
			if err == syscall.ENOENT {
				return nil, fs.ErrNotExist
			}
			return nil, err //fmt.Errorf("%w: %s", fmt.Errorf("%w", fs.ErrInvalid), err.Error())
		}
		return nil, nil
	}

	if flag == os.O_TRUNC {
		err := os.Remove(realPath)
		if err == nil {
			return nil, nil
		}
		if os.IsNotExist(err) {
			return nil, err
		}
		if err == syscall.ENOENT {
			return nil, fs.ErrNotExist
		}
		if err != syscall.ENOTEMPTY {
			return nil, &fs.PathError{
				Err:  fmt.Errorf("%w: directory not empty: %s", fs.ErrInvalid, err.Error()),
				Path: name,
			}
		}
		/*if err != nil {
			return nil, fmt.Errorf("%w: %s", fmt.Errorf("%w", fs.ErrInvalid), err.Error())
		}*/
		return nil, err
	}

	return os.OpenFile(realPath, flag, perm)
}

// DirWriteFS ...
func DirWriteFS(dir string) writefs.WriteFS {
	return osWriteFS{
		FS:   os.DirFS(dir),
		root: dir,
	}
}
