package writefs

import (
	"fmt"
	"io/fs"
	"os"
)

// RemoveFS is the interface implemented by a file system
// that provides an optimized implementation of Remove.
type RemoveFS interface {
	fs.FS
	Remove(name string, perm fs.FileMode) error
}

// Remove creates a new directory with the specified name
// and permission bits.
// If there is an error, it will be of type *PathError.
// If fsys implements RemoveFS, Remove calls fsys.Remove.
// Otherwise Remove calls fsys.OpenFile with a fs.O_DIR
// argument and fs.O_CREATE
func Remove(fsys fs.FS, name string, perm fs.FileMode) error {
	if fsys, ok := fsys.(RemoveFS); ok {
		return fsys.Remove(name, perm)
	}

	if fsys, ok := fsys.(WriteFS); ok {
		_, err := fsys.OpenFile(name, os.O_TRUNC, 0)
		return err
	}

	return fmt.Errorf("fsys does not support removal of files")
}
