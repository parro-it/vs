package mountedfs

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/parro-it/vs/memfs"
	"github.com/parro-it/vs/writefstest"
	"github.com/stretchr/testify/assert"
)

func TestMountedFS(t *testing.T) {
	data := []byte{0xca, 0xfe, 0xba, 0xbe}
	memfs1 := fstest.MapFS{
		"adir/afile": &fstest.MapFile{Data: data},
	}

	memfs2 := fstest.MapFS{
		"adir2/afile2": &fstest.MapFile{Data: data},
	}

	dir1 := memfs.New()
	dirempty := memfs.New()

	mfs := MountedFS{
		"c":        memfs1,
		"d":        memfs2,
		"dir1":     dir1,
		"dirempty": dirempty,
	}

	t.Run("pass writefstest.TestFS", writefstest.TestFS(mfs))

	t.Run("read files from multiple fs", func(t *testing.T) {
		buf, err := fs.ReadFile(mfs, "c/adir/afile")
		buf2, err2 := fs.ReadFile(mfs, "d/adir2/afile2")
		assert.NoError(t, err)
		assert.NoError(t, err2)
		assert.Equal(t, data, buf)
		assert.Equal(t, data, buf2)
	})

	t.Run("fs roots preserve their path", func(t *testing.T) {
		info, err := fs.Stat(mfs, "c")
		assert.NoError(t, err)
		assert.Equal(t, "c", info.Name())
	})

	t.Run("root contains all fs", func(t *testing.T) {
		info, err := fs.Stat(mfs, ".")
		assert.NoError(t, err)
		assert.Equal(t, ".", info.Name())
		assert.True(t, info.IsDir())
		assert.Equal(t, fs.ModeDir, info.Mode())
	})

	t.Run("unknown fs", func(t *testing.T) {
		buf, err := fs.ReadFile(mfs, "f/adir/afile")

		assert.Error(t, err)
		assert.Equal(t, "file does not exist: fs not found: f", err.Error())
		assert.Nil(t, buf)
	})

}
