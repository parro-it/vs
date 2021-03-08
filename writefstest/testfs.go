package writefstest

import (
	"errors"
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"

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

		dirNotExists := func(t *testing.T, dir string) {
			info, err := fs.Stat(fsys, dir)
			assert.True(t, os.IsNotExist(err))
			assert.Nil(t, info)
		}

		checkDirCreated := func(t *testing.T, dir string) {
			dirNotExists(t, dir)

			err := writefs.Remove(fsys, dir)
			assert.True(t, err == nil || os.IsNotExist(err))

			f, err := fsys.OpenFile(dir, os.O_CREATE, fs.FileMode(0755)|fs.ModeDir)
			assert.NoError(t, err)
			assert.Nil(t, f)
			dirExists(t, dir)
		}

		t.Run("creates directories with OpenFile - nested and not recursively", func(t *testing.T) {
			// nested dir return error
			f, err := fsys.OpenFile("adir/nested", os.O_CREATE, fs.FileMode(0755)|fs.ModeDir)
			assert.Error(t, err)
			assert.Nil(t, f)
			assert.True(t, errors.Is(err, fs.ErrNotExist))

			checkDirCreated(t, "adir")
			checkDirCreated(t, "adir/nested")
		})

		t.Run("OpenFile return *PathError on bad paths", func(t *testing.T) {
			checkBadPath(t, "afilename", "OpenFile", func(name string) error {
				_, err := fsys.OpenFile(name, 0, 0)
				return err
			})
		})

		checkDirRemoved := func(t *testing.T, dir string) {
			err := writefs.MkDir(fsys, dir, fs.FileMode(0755))
			assert.True(t, err == nil || os.IsExist(err))
			dirExists(t, dir)

			f, err := fsys.OpenFile(dir, os.O_TRUNC, 0)
			assert.NoError(t, err)
			assert.Nil(t, f)
			dirNotExists(t, dir)
		}

		t.Run("remove files with OpenFile", func(t *testing.T) {})

		t.Run("remove directories with OpenFile - nested and not recursively", func(t *testing.T) {
			// non empty dir return error
			f, err := fsys.OpenFile("adir", os.O_TRUNC, 0)
			assert.Error(t, err)
			assert.Nil(t, f)
			assert.True(t, errors.Is(err, fs.ErrInvalid))
			//assert.True(t, errors.Is(err, &fs.PathError{}))

			checkDirRemoved(t, "adir/nested")
			checkDirRemoved(t, "adir")
		})
		t.Run("create and write on new files", func(t *testing.T) {
			t.Run("set modtime to now", func(t *testing.T) {})
			t.Run("set content", func(t *testing.T) {})
		})
		t.Run("write on existing files", func(t *testing.T) {
			t.Run("updates modtime", func(t *testing.T) {})
			t.Run("update content", func(t *testing.T) {})
		})
		t.Run("write on existing files truncating", func(t *testing.T) {
			t.Run("updates modtime", func(t *testing.T) {})
			t.Run("set content", func(t *testing.T) {})
		})

		t.Run("appending to existing files", func(t *testing.T) {
			t.Run("updates modtime", func(t *testing.T) {})
			t.Run("update content", func(t *testing.T) {})
		})

		t.Run("opening non existing files", func(t *testing.T) {})
		t.Run("opening read-only files for write", func(t *testing.T) {})
		t.Run("", func(t *testing.T) {})
		t.Run("", func(t *testing.T) {})
		t.Run("", func(t *testing.T) {})
		t.Run("", func(t *testing.T) {})
		t.Run("", func(t *testing.T) {})
		t.Run("", func(t *testing.T) {})
		t.Run("", func(t *testing.T) {})
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
