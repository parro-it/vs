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
	/*
		t.Run("MapWriteFS", func(t *testing.T) {
			/ *		data := []byte{0xca, 0xfe, 0xba, 0xbe}

					fsw := MapWriteFS{fstest.MapFS{
						"adir/afile":   &fstest.MapFile{Data: data},
						"adir2/afile2": &fstest.MapFile{Data: data},
					}}
			* /
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
			/*
				t.Run("creates directories with OpenFile", func(t *testing.T) {

					fsys := MapWriteFS{fstest.MapFS{}}
					f, err := fsys.OpenFile("tmp/adir", os.O_CREATE, fs.FileMode(0755)|fs.ModeDir)
					assert.NoError(t, err)
					assert.Nil(t, f)
					info, err := fs.Stat(fsys, "tmp/adir")
					assert.NoError(t, err)
					assert.True(t, info.IsDir())

				})

				t.Run("remove directories with OpenFile", func(t *testing.T) {

					fsys := MapWriteFS{fstest.MapFS{
						"adir/afile": &fstest.MapFile{Mode: fs.ModeDir},
					}}

					info, err := fs.Stat(fsys, "adir/afile")
					assert.NoError(t, err)
					assert.True(t, info.IsDir())

					f, err := fsys.OpenFile("adir/afile", os.O_TRUNC, 0)
					assert.NoError(t, err)
					assert.Nil(t, f)

					info, err = fs.Stat(fsys, "adir/afile")
					assert.Error(t, err)
					assert.True(t, os.IsNotExist(err))
					assert.Nil(t, info)

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
			* /
		})
	*/
}
