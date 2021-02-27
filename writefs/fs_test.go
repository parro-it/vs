package writefs

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func fixtures() WriteFS {
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

func TestWriteFS(t *testing.T) {

	t.Run("defaults to Open", func(t *testing.T) {
		data := []byte{0xca, 0xfe, 0xba, 0xbe}

		memfs2 := fstest.MapFS{
			"adir2/afile2": &fstest.MapFile{Data: data},
		}

		f, err := OpenFile(memfs2, "adir2/afile2", os.O_RDONLY, fs.FileMode(0))
		defer f.Close()
		assert.NoError(t, err)
		buf := make([]byte, len(data))
		n, err := f.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, n, len(data))
		assert.Equal(t, data, buf)
	})

	t.Run("osWriteFS", func(t *testing.T) {
		fsw := fixtures()

		t.Run("writes files", func(t *testing.T) {
			os.Remove(fixtureFile("afile"))
			f, err := OpenFile(fsw, "afile", os.O_WRONLY|os.O_CREATE, fs.FileMode(0664))
			defer f.Close()

			buf := []byte("ciao")
			n, err := f.Write(buf)
			assert.NoError(t, err)
			assert.Equal(t, n, len(buf))
			actual, err := ioutil.ReadFile(fixtureFile("afile"))
			assert.NoError(t, err)
			assert.Equal(t, buf, actual)
			assert.NoError(t, os.Remove(fixtureFile("afile")))
		})

		t.Run("pass TestFS", func(t *testing.T) {
			fstest.TestFS(fsw, "dir1/file2", "file1")

			f, err := OpenFile(fsw, ".", os.O_RDONLY, fs.FileMode(0664))
			defer f.Close()
			assert.NoError(t, err)
			dir, ok := f.(fs.ReadDirFile)
			files, err := dir.ReadDir(-1)
			assert.True(t, ok)
			assert.Equal(t, 2, len(files))
			assert.Equal(t, true, files[1].Type().IsDir())
			assert.Equal(t, "dir1", files[1].Name())
			assert.Equal(t, false, files[0].Type().IsDir())
			assert.Equal(t, "file1", files[0].Name())

		})

	})

	t.Run("MapWriteFS", func(t *testing.T) {
		data := []byte{0xca, 0xfe, 0xba, 0xbe}

		fsw := MapWriteFS{fstest.MapFS{
			"adir/afile":   &fstest.MapFile{Data: data},
			"adir2/afile2": &fstest.MapFile{Data: data},
		}}

		t.Run("writes files", func(t *testing.T) {
			f, err := OpenFile(fsw, "adir/afile", os.O_WRONLY|os.O_TRUNC, fs.FileMode(0664))
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
