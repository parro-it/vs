package crudfs

/*
import (
	"io/fs"

	"github.com/parro-it/vs/writefs"
)

// RemoveWriteFS ...
type RemoveWriteFS interface {
	writefs.WriteFS
	Remove(name string) error
}

type osRemoveWriteFS struct {
	writefs.WriteFS
	root string
}

// MkDir implements RemoveWriteFS.MkDir
func (fsys osRemoveWriteFS) MkDir(name string) error {
	return nil
}

// Remove implements RemoveWriteFS.Remove
func (fsys osRemoveWriteFS) Remove(name string) error {

	return nil
}

// MapRemoveWriteFS ...
type MapRemoveWriteFS struct {
	writefs.MapWriteFS
}

// Remove implements RemoveWriteFS.Remove
func (fsys MapRemoveWriteFS) Remove(name string) error {
	return nil
}

// Remove ...
func Remove(fsys fs.FS, name string) error {
	/*
		if fs, ok := fsys.(WriteFS); ok {
				return fs.OpenFile(name, flag, perm)
			}

			if flag == os.O_RDONLY {
				file, err := fsys.Open(name)
				if err != nil {
					return nil, err
				}
				return readOnlyWriteFile{
					File: file,
				}, nil
			}

			return nil, fmt.Errorf("file system does not support write: %w", fs.ErrInvalid)
	* /
	return nil
}
*/
