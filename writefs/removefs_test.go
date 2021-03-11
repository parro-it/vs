package writefs

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestRemoveFS(t *testing.T) {
	roFS := fstest.MapFS{}
	writefs := testWriteFS{roFS, nil}
	removefs := &testRemoveFS{writefs, ""}

	t.Run("Remove calls Open with O_TRUNC for non RemoveFS but WriteFs instances", func(t *testing.T) {
		writefs.expectedErr = errors.New("adir2/afile2")
		err := Remove(writefs, "adir2/afile2")
		assert.Error(t, err)
		writefs.expectedErr = nil
		assert.Equal(t, "adir2/afile2", err.Error())
	})

	t.Run("Remove calls fsys.Remove for RemoveFS instances", func(t *testing.T) {
		removefs.removed = ""
		err := Remove(removefs, "adir2/afile2")
		assert.NoError(t, err)
		assert.Equal(t, "adir2/afile2", removefs.removed)
		removefs.removed = ""
	})

	t.Run("Remove return error for read onyl fs.FS", func(t *testing.T) {
		err := Remove(roFS, "adir2")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, fs.ErrInvalid))
		assert.Equal(t, "invalid argument: fsys does not support removal of files", err.Error())

	})
}
