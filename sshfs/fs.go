package sshfs

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/parro-it/vs/writefs"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SSHFS ...
type SSHFS struct {
	client        *sftp.Client
	ownedSSHCient *ssh.Client
	root          string
}

// OpenFile implements writefs.WriteFS
func (fsys *SSHFS) OpenFile(name string, flag int, perm fs.FileMode) (writefs.FileWriter, error) {
	return nil, nil
}

// Stat implements fs.StatFS
func (fsys *SSHFS) Stat(name string) (fs.FileInfo, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{}
	}
	return fsys.client.Stat(fsys.resolvePath(name))
}

// ReadFile implements fs.ReadFileFS
func (fsys *SSHFS) ReadFile(name string) ([]byte, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{}
	}

	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

/*
// Sub implements fs.SubFS
func (fsys *SSHFS) Sub(dir string) (fs.FS, error) {

	return nil, nil
}
*/

type fileWrapper struct {
	*sftp.File
	resolvedPath string
	fsys         *SSHFS
	dirContent   []fs.DirEntry
	cursor       int
}

type rootedDirEntry struct {
	os.FileInfo
}

func (f *fileWrapper) Stat() (fs.FileInfo, error) {
	info, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return rootedDirEntry{info}, nil

}

func (f *fileWrapper) ReadDir(n int) ([]fs.DirEntry, error) {
	if f.dirContent == nil {
		files, err := f.fsys.client.ReadDir(f.resolvedPath)
		if err != nil {
			return nil, err
		}
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})
		f.dirContent = make([]fs.DirEntry, len(files))
		for idx, file := range files {
			f.dirContent[idx] = fileInfoDirEntry{file}
		}
	}
	if n == -1 {
		if f.cursor > 0 {
			return []fs.DirEntry{}, nil
		}
		f.cursor += len(f.dirContent)
		return f.dirContent, nil
	}
	if f.cursor >= len(f.dirContent) {
		return []fs.DirEntry{}, io.EOF
	}
	last := f.cursor + n
	var err error
	if last > len(f.dirContent) {
		last = len(f.dirContent)
		err = io.EOF
	}

	res := f.dirContent[f.cursor:last]
	f.cursor += n

	return res, err
}

// Open implements fs.FS
func (fsys *SSHFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{}
	}

	f, err := fsys.client.Open(fsys.resolvePath(name))
	if err != nil {
		return nil, err
	}

	wrapper := fileWrapper{
		File:         f,
		resolvedPath: fsys.resolvePath(name),
		fsys:         fsys,
	}

	return &wrapper, nil

}

func (fsys *SSHFS) resolvePath(name string) string {
	path := name
	if path == "." {
		path = fsys.root
	} else {
		path = filepath.Join(fsys.root, name)
	}
	return path
}

// ReadDir implements fs.ReadDirFS
func (fsys *SSHFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{}
	}

	files, err := fsys.client.ReadDir(fsys.resolvePath(name))
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	ret := make([]fs.DirEntry, len(files))
	for idx, f := range files {
		ret[idx] = fileInfoDirEntry{f}
	}

	return ret, nil
}

type fileInfoDirEntry struct {
	os.FileInfo
}

// Type returns the type bits for the entry.
// The type bits are a subset of the usual FileMode bits, those returned by the FileMode.Type method.
func (e fileInfoDirEntry) Type() fs.FileMode {
	return e.Mode().Type()
}

// Info returns the FileInfo for the file or subdirectory described by the entry.
// The returned FileInfo may be from the time of the original directory read
// or from the time of the call to Info. If the file has been removed or renamed
// since the directory read, Info may return an error satisfying errors.Is(err, ErrNotExist).
// If the entry denotes a symbolic link, Info reports the information about the link itself,
// not the link's target.
func (e fileInfoDirEntry) Info() (fs.FileInfo, error) {
	return e.FileInfo, nil
}
