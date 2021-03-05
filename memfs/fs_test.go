package memfs

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"testing/fstest"

	"github.com/parro-it/vs/writefs"
	"github.com/stretchr/testify/assert"
)

func fixtureFile(name string) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot retrieve the source file path")
	}

	file = filepath.Dir(filepath.Dir(file))

	return path.Join(file, "fixtures", name)
}

func TestOsFS(t *testing.T) {

	t.Run("MapWriteFS", func(t *testing.T) {
		data := []byte{0xca, 0xfe, 0xba, 0xbe}

		fsw := MapWriteFS{fstest.MapFS{
			"adir/afile":   &fstest.MapFile{Data: data},
			"adir2/afile2": &fstest.MapFile{Data: data},
		}}

		t.Run("WriteFile", func(t *testing.T) {
			fsw := MapWriteFS{fstest.MapFS{}}

			buf := []byte("pippero")
			n, err := writefs.WriteFile(fsw, "file.pip", buf)
			assert.NoError(t, err)
			assert.Equal(t, len(buf), n)

			actual, err := fs.ReadFile(fsw, "file.pip")
			assert.NoError(t, err)
			assert.Equal(t, buf, actual)

		})

		t.Run("writes files", func(t *testing.T) {
			f, err := writefs.OpenFile(fsw, "adir/afile", os.O_WRONLY|os.O_TRUNC, fs.FileMode(0664))
			assert.NoError(t, err)

			defer f.Close()

			buf := []byte("ciao")
			n, err := f.Write(buf)

			assert.NoError(t, err)
			assert.Equal(t, n, len(buf))

			actual, err := fs.ReadFile(fsw, "adir/afile")
			assert.NoError(t, err)
			assert.Equal(t, buf, actual)
		})

		t.Run("pass TestFS", func(t *testing.T) {
			fstest.TestFS(fsw, "adir/afile", "adir2/afile2")
		})

	})

}
