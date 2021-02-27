package mountedfs

import (
	"fmt"
	"io/fs"
	"strings"
	"syscall"
	"testing/fstest"

	"github.com/parro-it/vs/writefs"
)

// MountedFS is a `fs.FS` implementation that
// allows to group together multiple named filesystems
// and access them as if they are mounted under a
// virtual root directory, whose entries are names
// as the filesystems theirself.
//
//
//	mfs := MountedFS{
//		"mem1": fstest.MapFS{
//			"adir/afile": &fstest.MapFile{Data: data},
//		},
//		"mem2": fstest.MapFS{
//			"adir2/afile2": &fstest.MapFile{Data: data},
//		},
//	}
//
//	mfs.Open("mem1/adir/afile")			// returns adir/afile from mem1 fs
//	mfs.Open("mem2/adir2/afile2")		// returns adir2/afile2 from mem2 fs
//	mfs.Open("mem2") 					// returns root of mem2
//	mfs.Open(".") 						// returns a virtual ReadDirFile containing
//										// an entry for mem1 and one for mem2
//
// MountedFS implements:
// * fs.ReadFileFS
// * fs.StatFS
// * fs.SubFS
// * writefs.WriteFS
type MountedFS map[string]fs.FS

// Stat implements fs.StatFS
func (f MountedFS) Stat(name string) (fs.FileInfo, error) {
	if name == "." {
		// when requested file is the root of
		// its fs, return a mem FileInfo that
		// adjust its file name.
		return newMemDirInfo("."), nil
	}
	rpath := f.pickRemotePath(name)
	if rpath.Error != nil {
		return nil, rpath.Error
	}
	if rpath.Path == "." {
		// when requested file is the root of
		// its fs, return a mem FileInfo that
		// adjust its file name.
		return newMemDirInfo(rpath.FsName), nil
	}

	return fs.Stat(rpath.Fs, rpath.Path)
}

// ReadFile implements fs.ReadFileFS
func (f MountedFS) ReadFile(name string) ([]byte, error) {
	if name == "." {
		return nil, syscall.EISDIR
	}
	rpath := f.pickRemotePath(name)
	if rpath.Error != nil {
		return nil, rpath.Error
	}
	if rpath.Path == "." {
		return nil, syscall.EISDIR
	}

	return fs.ReadFile(rpath.Fs, rpath.Path)
}

// Sub implements fs.SubFS
func (f MountedFS) Sub(dir string) (fs.FS, error) {
	if dir == "." {
		return f, nil
	}
	rpath := f.pickRemotePath(dir)
	if rpath.Error != nil {
		return nil, rpath.Error
	}

	if rpath.Path == "." {
		return rpath.Fs, nil
	}

	return fs.Sub(rpath.Fs, rpath.Path)
}

// OpenFile implements writefs.WriteFS
func (f MountedFS) OpenFile(name string, flag int, perm fs.FileMode) (writefs.WriteFile, error) {
	return nil, nil
}

// Mount add a child file system, using `name`
// argument as it's mount name.
func (f MountedFS) Mount(name string, fs fs.FS) {
	f[name] = fs
}

// Open opens the named file.
func (f MountedFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{}
	}

	if name == "." {
		// special case for root dir, returns list
		// of mounted fs
		return newVirtualRootDir(f), nil
	}

	rpath := f.pickRemotePath(name)
	if rpath.Error != nil {
		return nil, rpath.Error
	}
	file, err := rpath.Fs.Open(rpath.Path)
	if rpath.Path == "." {
		// when requested file is the root of
		// its fs, embed it in a virtualDir to handle
		// subdir correctly
		file = &virtualDir{file.(fs.ReadDirFile), rpath.FsName}
	}
	return file, err

}

// a virtual dir wraps
// a fs.ReadDirFile instance
// to add the name of its FS
// as first segment of its name.
type virtualDir struct {
	fs.ReadDirFile
	name string
}

func (f virtualDir) Stat() (fs.FileInfo, error) {
	//info, err := f.ReadDirFile.Stat()
	return newMemDirInfo(f.name), nil
}

func newMemDirFile(name string) fs.File {
	tmp := fstest.MapFS{}
	tmp[name] = &fstest.MapFile{
		Data: []byte{},
		Mode: fs.ModeDir,
	}
	file, err := tmp.Open(name)
	if err != nil {
		panic(err)
	}
	return file
}

func newMemDirEntry(name string) fs.DirEntry {
	tmp := fstest.MapFS{}
	tmp[name] = &fstest.MapFile{
		Data: []byte{},
		Mode: fs.ModeDir,
	}
	files, err := tmp.ReadDir(".")
	if err != nil {
		panic(err)
	}
	return files[0]
}

func newMemDirInfo(name string) fs.FileInfo {
	tmpFile := newMemDirFile(name)
	info, err := tmpFile.Stat()
	if err != nil {
		panic(err)
	}
	return info
}

func newVirtualRootDir(f MountedFS) fs.ReadDirFile {

	dir := fstest.MapFS{}
	for name := range f {
		dir[name] = &fstest.MapFile{
			Data: []byte{},
			Mode: fs.ModeDir,
		}
	}
	file, err := dir.Open(".")
	if err != nil {
		panic(err)
	}
	readDir, ok := file.(fs.ReadDirFile)
	if !ok {
		panic(err)
	}

	return readDir
}

type remotePath struct {
	Fs     fs.FS
	FsName string
	Path   string
	Error  error
}

func (f MountedFS) pickRemotePath(name string) remotePath {
	res := remotePath{}

	parts := strings.Split(name, "/")
	res.FsName = parts[0] // first segment identify the fs
	if len(parts) == 1 {
		res.Path = "."
	} else {
		res.Path = strings.Join(parts[1:], "/") // segments after first identify
		//				   					// the real path within the fs

	}
	var ok bool
	res.Fs, ok = f[res.FsName]
	if !ok {
		res.Error = fmt.Errorf("fs not found: %s", res.FsName)
	}

	return res
}
