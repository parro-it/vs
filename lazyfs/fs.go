package genericfs

import (
	"io/fs"

	"github.com/parro-it/vs/writefs"
)

type fsT struct {
	Factory func() (fs.FS, error)
	wrapped fs.FS
}

// New ...
func New(factory func() (fs.FS, error)) fs.FS {
	return &fsT{Factory: factory}
}

var (
	_ fs.StatFS     = &fsT{}
	_ fs.ReadFileFS = &fsT{}
	_ fs.SubFS      = &fsT{}
	_ fs.ReadDirFS  = &fsT{}
	_ fs.GlobFS     = &fsT{}

	_ writefs.WriteFS  = &fsT{}
	_ writefs.RemoveFS = &fsT{}
	_ writefs.MkDirFS  = &fsT{}
)

func (fsys *fsT) init() error {
	if fsys.wrapped != nil {
		return nil
	}

	f, err := fsys.Factory()
	if err != nil {
		return err
	}
	fsys.wrapped = f
	return nil
}

// MkDir implements writefs.MkDirFS
func (fsys *fsT) MkDir(name string, perm fs.FileMode) error {
	if err := fsys.init(); err != nil {
		return err
	}
	return writefs.MkDir(fsys.wrapped, name, perm)
}

// Remove implements writefs.RemoveFS
func (fsys *fsT) Remove(name string) error {
	if err := fsys.init(); err != nil {
		return err
	}
	return writefs.Remove(fsys.wrapped, name)
}

// OpenFile implements writefs.WriteFS
func (fsys *fsT) OpenFile(name string, flag int, perm fs.FileMode) (writefs.FileWriter, error) {
	if err := fsys.init(); err != nil {
		return nil, err
	}
	return writefs.OpenFile(fsys.wrapped, name, flag, perm)
}

// Stat implements fs.StatFS
func (fsys *fsT) Stat(name string) (fs.FileInfo, error) {
	if err := fsys.init(); err != nil {
		return nil, err
	}
	return fs.Stat(fsys.wrapped, name)
}

// ReadFile implements fs.ReadFileFS
func (fsys *fsT) ReadFile(name string) ([]byte, error) {
	if err := fsys.init(); err != nil {
		return nil, err
	}
	return fs.ReadFile(fsys.wrapped, name)
}

// Sub implements fs.SubFS
func (fsys *fsT) Sub(dir string) (fs.FS, error) {
	if err := fsys.init(); err != nil {
		return nil, err
	}
	return fs.Sub(fsys.wrapped, dir)
}

// Open implements fs.FS
func (fsys *fsT) Open(name string) (fs.File, error) {
	if err := fsys.init(); err != nil {
		return nil, err
	}
	return fsys.wrapped.Open(name)
}

// ReadDir implements fs.ReadDirFS
func (fsys *fsT) ReadDir(name string) ([]fs.DirEntry, error) {
	if err := fsys.init(); err != nil {
		return nil, err
	}
	return fs.ReadDir(fsys.wrapped, name)
}

// Glob implements fs.GlobFS
func (fsys *fsT) Glob(pattern string) ([]string, error) {
	if err := fsys.init(); err != nil {
		return nil, err
	}
	return fs.Glob(fsys.wrapped, pattern)
}
