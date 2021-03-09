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
	/*
		fsw := fixtures()

		t.Run("writes files", func(t *testing.T) {
			os.Remove(fixtureFile("afile"))
			f, err := writefs.OpenFile(fsw, "afile", os.O_WRONLY|os.O_CREATE, fs.FileMode(0664))
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

		t.Run("creates directories with OpenFile", func(t *testing.T) {
			os.RemoveAll("/tmp/adir")
			fsys := DirWriteFS("/tmp")
			f, err := fsys.OpenFile("adir", os.O_CREATE, fs.FileMode(0755)|fs.ModeDir)
			assert.NoError(t, err)
			assert.Nil(t, f)
			info, err := os.Stat("/tmp/adir")
			assert.NoError(t, err)
			assert.True(t, info.IsDir())
			err = os.RemoveAll("/tmp/adir")
			assert.NoError(t, err)
		})

		t.Run("remove directories with OpenFile", func(t *testing.T) {
			os.MkdirAll("/tmp/adir", os.FileMode(0755))
			info, err := os.Stat("/tmp/adir")
			assert.NoError(t, err)
			assert.True(t, info.IsDir())

			fsys := DirWriteFS("/tmp")
			f, err := fsys.OpenFile("adir", os.O_TRUNC, 0)
			assert.NoError(t, err)
			assert.Nil(t, f)

			info, err = os.Stat("/tmp/adir")
			assert.Error(t, err)
			assert.True(t, os.IsNotExist(err))
			assert.Nil(t, info)

			os.RemoveAll("/tmp/adir")

		})

		t.Run("pass TestFS", func(t *testing.T) {
			fstest.TestFS(fsw, "dir1/file2", "file1")

			f, err := writefs.OpenFile(fsw, ".", os.O_RDONLY, fs.FileMode(0664))
			defer f.Close()
			assert.NoError(t, err)
			dir, ok := f.(fs.ReadDirFile)
			files, err := dir.ReadDir(-1)
			assert.True(t, ok)
			assert.Equal(t, 4, len(files))
			assert.Equal(t, true, files[0].Type().IsDir())
			assert.Equal(t, "fakehost", files[0].Name())
			assert.Equal(t, false, files[1].Type().IsDir())
			assert.Equal(t, "file1", files[1].Name())
			assert.Equal(t, true, files[2].Type().IsDir())
			assert.Equal(t, "dir1", files[2].Name())
			assert.Equal(t, false, files[3].Type().IsDir())
			assert.Equal(t, "anyfile.txt", files[3].Name())

		})
	*/
}
