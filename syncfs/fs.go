package syncfs

import (
	"io/fs"
	"sync"

	"github.com/parro-it/vs/writefs"
)

type fsT struct {
	lock   sync.Mutex
	wrapfs fs.FS
}

// New ...
func New(fsys fs.FS) fs.FS {
	return &fsT{
		lock:   sync.Mutex{},
		wrapfs: fsys,
	}
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

// MkDir implements writefs.MkDirFS
func (fsys *fsT) MkDir(name string, perm fs.FileMode) error {
	fsys.lock.Lock()
	defer fsys.lock.Unlock()

	return writefs.MkDir(fsys.wrapfs, name, perm)
}

// Remove implements writefs.RemoveFS
func (fsys *fsT) Remove(name string) error {
	fsys.lock.Lock()
	defer fsys.lock.Unlock()
	return writefs.Remove(fsys.wrapfs, name)
}

// OpenFile implements writefs.WriteFS
func (fsys *fsT) OpenFile(name string, flag int, perm fs.FileMode) (writefs.FileWriter, error) {
	fsys.lock.Lock()
	defer fsys.lock.Unlock()
	return writefs.OpenFile(fsys.wrapfs, name, flag, perm)
}

// Stat implements fs.StatFS
func (fsys *fsT) Stat(name string) (fs.FileInfo, error) {
	fsys.lock.Lock()
	defer fsys.lock.Unlock()
	return fs.Stat(fsys.wrapfs, name)
}

// ReadFile implements fs.ReadFileFS
func (fsys *fsT) ReadFile(name string) ([]byte, error) {
	fsys.lock.Lock()
	defer fsys.lock.Unlock()
	return fs.ReadFile(fsys.wrapfs, name)
}

// Sub implements fs.SubFS
func (fsys *fsT) Sub(dir string) (fs.FS, error) {
	fsys.lock.Lock()
	defer fsys.lock.Unlock()
	return fs.Sub(fsys.wrapfs, dir)
}

// Open implements fs.FS
func (fsys *fsT) Open(name string) (fs.File, error) {
	fsys.lock.Lock()
	defer fsys.lock.Unlock()
	return fsys.wrapfs.Open(name)
}

// ReadDir implements fs.ReadDirFS
func (fsys *fsT) ReadDir(name string) ([]fs.DirEntry, error) {
	fsys.lock.Lock()
	defer fsys.lock.Unlock()
	return fs.ReadDir(fsys.wrapfs, name)
}

// Glob implements fs.GlobFS
func (fsys *fsT) Glob(pattern string) ([]string, error) {
	fsys.lock.Lock()
	defer fsys.lock.Unlock()
	return fs.Glob(fsys.wrapfs, pattern)
}
