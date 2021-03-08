package genericfs

import (
	"io/fs"

	"github.com/parro-it/vs/writefs"
)

// FS ...
type FS struct{}

var (
	_ fs.StatFS        = FS{}
	_ fs.ReadFileFS    = FS{}
	_ fs.SubFS         = FS{}
	_ writefs.WriteFS  = FS{}
	_ writefs.RemoveFS = FS{}
	_ writefs.MkDirFS  = FS{}
	_ fs.ReadDirFS     = FS{}
	_ fs.GlobFS        = FS{}
)

// MkDir implements writefs.MkDirFS
func (fsys FS) MkDir(name string, perm fs.FileMode) error {
	return nil
}

// Remove implements writefs.RemoveFS
func (fsys FS) Remove(name string) error {
	return nil
}

// OpenFile implements writefs.WriteFS
func (fsys FS) OpenFile(name string, flag int, perm fs.FileMode) (writefs.FileWriter, error) {
	return nil, nil
}

// Stat implements fs.StatFS
func (fsys FS) Stat(name string) (fs.FileInfo, error) {
	return nil, nil
}

// ReadFile implements fs.ReadFileFS
func (fsys FS) ReadFile(name string) ([]byte, error) {
	return nil, nil
}

// Sub implements fs.SubFS
func (fsys FS) Sub(dir string) (fs.FS, error) {
	return nil, nil
}

// Open implements fs.FS
func (fsys FS) Open(name string) (fs.File, error) {
	return nil, nil
}

// ReadDir implements fs.ReadDirFS
func (fsys FS) ReadDir(name string) ([]fs.DirEntry, error) {
	return nil, nil
}

// Glob implements fs.GlobFS
func (fsys FS) Glob(pattern string) ([]string, error) {
	return nil, nil
}
