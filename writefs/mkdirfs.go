package writefs

import (
	"fmt"
	"io/fs"
	"os"
)

// MkDirFS is the interface implemented by a file system
// that provides an optimized implementation of MkDir.
type MkDirFS interface {
	fs.FS
	MkDir(name string, perm fs.FileMode) error
}

// MkDir creates a new directory with the specified name
// and permission bits.
// If there is an error, it will be of type *PathError.
// If fsys implements MkDirFS, MkDir calls fsys.MkDir.
// Otherwise MkDir calls fsys.OpenFile with a fs.O_DIR
// argument and fs.O_CREATE
func MkDir(fsys fs.FS, name string, perm fs.FileMode) error {
	if fsys, ok := fsys.(MkDirFS); ok {
		return fsys.MkDir(name, perm)
	}

	if fsys, ok := fsys.(WriteFS); ok {
		f, err := fsys.OpenFile(name, os.O_CREATE, perm|fs.ModeDir)
		if f != nil {
			f.Close()
		}
		return err
	}

	return fmt.Errorf("%w: fsys does not support creation of directories", fs.ErrInvalid)
}
