package osfs

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/parro-it/vs/writefs"
	"github.com/parro-it/vs/writefstest"
)

func fixtures() writefs.WriteFS {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot retrieve the source file path")
	}

	file = filepath.Dir(filepath.Dir(file))

	return DirWriteFS(path.Join(file, "fixtures"))
}

func fixtureFile(name string) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot retrieve the source file path")
	}

	file = filepath.Dir(filepath.Dir(file))

	return path.Join(file, "fixtures", name)
}

func TestOsFS(t *testing.T) {
	fsys := DirWriteFS("/var/fixtures")
	t.Run("Pass writefstest.TestFS", writefstest.TestFS(fsys))

}
