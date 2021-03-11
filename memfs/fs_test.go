package memfs

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/parro-it/vs/writefstest"
)

func fixtureFile(name string) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot retrieve the source file path")
	}

	file = filepath.Dir(filepath.Dir(file))

	return path.Join(file, "fixtures", name)
}

func TestMemFS(t *testing.T) {
	fsys := New()
	t.Run("Pass writefstest.TestFS", writefstest.TestFS(fsys))

}
