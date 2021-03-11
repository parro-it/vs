package writefs

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestMkDirFS(t *testing.T) {
	roFS := fstest.MapFS{}
	writefs := testWriteFS{roFS, nil}
	mkdirfs := &testMkDirFS{writefs, ""}

	t.Run("MkDir calls Open with O_CREATE for non RemoveFS but WriteFs instances", func(t *testing.T) {
		writefs.expectedErr = errors.New("adir2/afile2")
		err := MkDir(writefs, "adir2/afile2", fs.FileMode(0))
		assert.Error(t, err)
		writefs.expectedErr = nil
		assert.Equal(t, "adir2/afile2", err.Error())
	})

	t.Run("MkDir calls fsys.MkDir for RemoveFS instances", func(t *testing.T) {
		mkdirfs.created = ""
		err := MkDir(mkdirfs, "adir2/afile2", fs.FileMode(0))
		assert.NoError(t, err)
		assert.Equal(t, "adir2/afile2", mkdirfs.created)
		mkdirfs.created = ""
	})

	t.Run("MkDir return error for read only fs.FS", func(t *testing.T) {
		err := MkDir(roFS, "adir2", fs.FileMode(0))
		assert.Error(t, err)
		assert.True(t, errors.Is(err, fs.ErrInvalid))
		assert.Equal(t, "invalid argument: fsys does not support creation of directories", err.Error())

	})

}
