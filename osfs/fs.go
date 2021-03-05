package osfs

import (
	"io/fs"
	"os"
	"path"

	"github.com/parro-it/vs/writefs"
)

type osWriteFS struct {
	fs.FS
	root string
}

// OpenFile ...
func (fsinst osWriteFS) OpenFile(name string, flag int, perm fs.FileMode) (writefs.FileWriter, error) {
	realPath := path.Join(fsinst.root, name)
	return os.OpenFile(realPath, flag, perm)
}

// DirWriteFS ...
func DirWriteFS(dir string) writefs.WriteFS {
	return osWriteFS{
		FS:   os.DirFS(dir),
		root: dir,
	}
}