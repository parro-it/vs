package writefstest

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/parro-it/vs/writefs"
	"github.com/stretchr/testify/assert"
)

// TestFS returns a function that test given writefs.WriteFS
// with common operation, like fstest.TestFS does for readonly FSs
func TestFS(fsys writefs.WriteFS) func(t *testing.T) {
	return func(t *testing.T) {
		dirs := []string{
			"dir1",
			"dir1/dirsub1",
			"dirempty",
		}
		files := []string{
			"dir1/file1",
			"dir1/file2",
			"dir1/dirsub1/file3",
		}

		t.Run("initialize testing FS", func(t *testing.T) {
			for _, dir := range dirs {
				err := writefs.MkDir(fsys, dir, fs.FileMode(0755))
				if !(err == nil || errors.Is(err, fs.ErrExist)) {
					fmt.Println(err, dir)
				}
				assert.True(t, err == nil || errors.Is(err, fs.ErrExist))

			}

			for _, file := range files {
				fake := []byte(file + " content\n")
				n, err := writefs.WriteFile(fsys, file, fake)
				assert.NoError(t, err)
				assert.Equal(t, len(fake), n)
			}
		})

		t.Run("pass TestFS", func(t *testing.T) {
			err := fstest.TestFS(fsys, append(files, dirs...)...)
			assert.NoError(t, err)
		})

		dirExists := func(t *testing.T, dir string) {
			info, err := fs.Stat(fsys, dir)
			assert.NoError(t, err)
			assert.True(t, info.IsDir())
		}

		fileExists := func(t *testing.T, dir string) {
			info, err := fs.Stat(fsys, dir)
			assert.NoError(t, err)
			assert.True(t, info.Mode().IsRegular())
		}

		fileNotExists := func(t *testing.T, dir string) {
			info, err := fs.Stat(fsys, dir)
			assert.True(t, os.IsNotExist(err))
			assert.Nil(t, info)
		}

		checkDirCreated := func(t *testing.T, dir string) {
			fileNotExists(t, dir)

			err := writefs.Remove(fsys, dir)
			assert.True(t, err == nil || os.IsNotExist(err))

			f, err := fsys.OpenFile(dir, os.O_CREATE, fs.FileMode(0755)|fs.ModeDir)
			assert.NoError(t, err)
			assert.Nil(t, f)
			dirExists(t, dir)
		}
		dirRemove := func(t *testing.T, dir string) {
			f, _ := fsys.OpenFile(dir, os.O_TRUNC, 0)
			assert.Nil(t, f)
			fileNotExists(t, dir)
		}
		checkDirRemoved := func(t *testing.T, dir string) {
			err := writefs.MkDir(fsys, dir, fs.FileMode(0755))
			assert.True(t, err == nil || os.IsExist(err))
			dirExists(t, dir)

			f, err := fsys.OpenFile(dir, os.O_TRUNC, 0)
			assert.NoError(t, err)
			assert.Nil(t, f)
			fileNotExists(t, dir)
		}
		t.Run("creates directories with OpenFile - nested and not recursively", func(t *testing.T) {
			dirRemove(t, "dir1/adir/nested")
			dirRemove(t, "dir1/adir")

			// nested dir return error
			f, err := fsys.OpenFile("dir1/adir/nested", os.O_CREATE, fs.FileMode(0755)|fs.ModeDir)
			assert.Error(t, err)
			assert.Nil(t, f)
			assert.True(t, errors.Is(err, fs.ErrNotExist))

			checkDirCreated(t, "dir1/adir")
			checkDirCreated(t, "dir1/adir/nested")
		})

		t.Run("OpenFile return *PathError on bad paths", func(t *testing.T) {
			checkBadPath(t, "afilename", "OpenFile", func(name string) error {
				_, err := fsys.OpenFile(name, 0, 0)
				return err
			})
		})

		t.Run("remove files with OpenFile", func(t *testing.T) {
			file := "dir1/somenewfile"
			_, err := writefs.WriteFile(fsys, file, []byte(file))
			assert.NoError(t, err)

			fileExists(t, file)

			f, err := fsys.OpenFile(file, os.O_TRUNC, 0)
			assert.NoError(t, err)
			assert.Nil(t, f)

			fileNotExists(t, file)
		})

		t.Run("remove directories with OpenFile - nested and not recursively", func(t *testing.T) {
			// non empty dir return error
			f, err := fsys.OpenFile("dir1/adir", os.O_TRUNC, 0)
			assert.Error(t, err)
			assert.Nil(t, f)
			assert.True(t, errors.Is(err, fs.ErrInvalid))
			//assert.True(t, errors.Is(err, &fs.PathError{}))

			checkDirRemoved(t, "dir1/adir/nested")
			checkDirRemoved(t, "dir1/adir")
		})
		t.Run("create and write on new files", func(t *testing.T) {
			file := "dir1/file1new"
			err := writefs.Remove(fsys, file)
			assert.True(t, err == nil || os.IsNotExist(err))
			fileNotExists(t, file)

			f, err := fsys.OpenFile(file, os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
			assert.NoError(t, err)
			assert.NotNil(t, f)
			buf := []byte("ciao\n")
			n, err := f.Write(buf)
			assert.NoError(t, err)
			assert.Equal(t, len(buf), n)
			err = f.Close()
			assert.NoError(t, err)

			t.Run("set modtime to now", func(t *testing.T) {
				info, err := fs.Stat(fsys, file)
				assert.NoError(t, err)
				assert.Less(t, time.Now().Sub(info.ModTime()), time.Second)
			})
			t.Run("set content", func(t *testing.T) {
				actual, err := fs.ReadFile(fsys, file)
				assert.NoError(t, err)
				assert.Equal(t, buf, actual)
			})
		})

		t.Run("write on existing files", func(t *testing.T) {
			file := "dir1/file1new"
			err := writefs.Remove(fsys, file)
			assert.True(t, err == nil || os.IsNotExist(err))
			_, err = writefs.WriteFile(fsys, file, []byte("ciao\n"))
			assert.NoError(t, err)

			fileExists(t, file)

			f, err := fsys.OpenFile(file, os.O_WRONLY, os.FileMode(0644))
			assert.NoError(t, err)
			assert.NotNil(t, f)
			buf := []byte("mi")
			n, err := f.Write(buf)
			assert.NoError(t, err)
			assert.Equal(t, len(buf), n)
			err = f.Close()
			assert.NoError(t, err)

			t.Run("updates modtime", func(t *testing.T) {
				info, err := fs.Stat(fsys, file)
				assert.NoError(t, err)
				assert.Less(t, time.Now().Sub(info.ModTime()), time.Second)
			})
			t.Run("update content", func(t *testing.T) {
				actual, err := fs.ReadFile(fsys, file)
				assert.NoError(t, err)
				assert.Equal(t, []byte("miao\n"), actual)
			})
		})
		t.Run("write on existing files truncating", func(t *testing.T) {
			file := "dir1/file1new"
			err := writefs.Remove(fsys, file)
			assert.True(t, err == nil || os.IsNotExist(err))
			_, err = writefs.WriteFile(fsys, file, []byte("ciao\n"))
			assert.NoError(t, err)

			fileExists(t, file)

			f, err := fsys.OpenFile(file, os.O_WRONLY|os.O_TRUNC, os.FileMode(0644))
			assert.NoError(t, err)
			assert.NotNil(t, f)
			buf := []byte("mi")
			n, err := f.Write(buf)
			assert.NoError(t, err)
			assert.Equal(t, len(buf), n)
			err = f.Close()
			assert.NoError(t, err)

			t.Run("updates modtime", func(t *testing.T) {
				info, err := fs.Stat(fsys, file)
				assert.NoError(t, err)
				assert.Less(t, time.Now().Sub(info.ModTime()), time.Second)
			})
			t.Run("set content", func(t *testing.T) {
				actual, err := fs.ReadFile(fsys, file)
				assert.NoError(t, err)
				assert.Equal(t, []byte("mi"), actual)
			})
		})

		t.Run("appending to existing files", func(t *testing.T) {
			file := "dir1/file1new"
			err := writefs.Remove(fsys, file)
			assert.True(t, err == nil || os.IsNotExist(err))
			_, err = writefs.WriteFile(fsys, file, []byte("ciao\n"))
			assert.NoError(t, err)

			fileExists(t, file)

			f, err := fsys.OpenFile(file, os.O_WRONLY|os.O_APPEND, os.FileMode(0644))
			assert.NoError(t, err)
			assert.NotNil(t, f)
			buf := []byte("mi")
			n, err := f.Write(buf)
			assert.NoError(t, err)
			assert.Equal(t, len(buf), n)
			err = f.Close()
			assert.NoError(t, err)

			t.Run("updates modtime", func(t *testing.T) {
				info, err := fs.Stat(fsys, file)
				assert.NoError(t, err)
				assert.Less(t, time.Now().Sub(info.ModTime()), time.Second)
			})
			t.Run("updates content", func(t *testing.T) {
				actual, err := fs.ReadFile(fsys, file)
				assert.NoError(t, err)
				assert.Equal(t, []byte("ciao\nmi"), actual)
			})
		})

		t.Run("opening non existing files", func(t *testing.T) {
			f, err := fsys.OpenFile("unkfile", os.O_WRONLY, os.FileMode(0644))
			assert.Error(t, err)
			assert.True(t, errors.Is(err, fs.ErrNotExist))
			assert.Nil(t, f)
		})
		/*
			t.Run("opening read-only files for write", func(t *testing.T) {
				f, err := fsys.OpenFile("/etc/passwd", os.O_WRONLY, os.FileMode(0644))
				assert.Error(t, err)
				assert.Nil(t, f)
			})
		*/
	}
}

// checkBadPath checks that various invalid forms of file's name cannot be opened using open.
func checkBadPath(t *testing.T, file string, desc string, open func(string) error) {
	bad := []string{
		"/" + file,
		file + "/.",
	}
	if file == "." {
		bad = append(bad, "/")
	}
	if i := strings.Index(file, "/"); i >= 0 {
		bad = append(bad,
			file[:i]+"//"+file[i+1:],
			file[:i]+"/./"+file[i+1:],
			file[:i]+`\`+file[i+1:],
			file[:i]+"/../"+file,
		)
	}
	if i := strings.LastIndex(file, "/"); i >= 0 {
		bad = append(bad,
			file[:i]+"//"+file[i+1:],
			file[:i]+"/./"+file[i+1:],
			file[:i]+`\`+file[i+1:],
			file+"/../"+file[i+1:],
		)
	}

	for _, b := range bad {
		assert.Error(t, open(b))
	}
}
