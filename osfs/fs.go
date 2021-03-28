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
			//return nil, fsErr(err, name)
			if os.IsExist(err) {
				return nil, fsErr(err, name)
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
		//return nil, fsErr(err, name)

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

		return nil, err

	}

	f, err := os.OpenFile(realPath, flag, perm)
	if err != nil {
		return nil, err //fsErr(err, name)
	}
	return f, nil
}

func fsErr(err error, name string) error {
	if os.IsExist(err) {
		return fmt.Errorf("%w: %s", fs.ErrExist, name)
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", fs.ErrNotExist, name)
	}
	if os.IsPermission(err) {
		return fmt.Errorf("%w: %s", fs.ErrPermission, name)
	}

	return err
}

// DirWriteFS ...
func DirWriteFS(dir string) writefs.WriteFS {
	return osWriteFS{
		FS:   os.DirFS(dir),
		root: dir,
	}
}
