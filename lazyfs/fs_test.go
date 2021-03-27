package genericfs

import (
	"io/fs"
	"testing"

	"github.com/parro-it/vs/memfs"
	"github.com/parro-it/vs/writefs"
	"github.com/parro-it/vs/writefstest"
)

func TestLazyFS(t *testing.T) {
	t.Run("pass writefstest.TestFS", func(t *testing.T) {
		fsys := New(func() (fs.FS, error) {
			return memfs.New(), nil
		})
		t.Run("Pass writefstest.TestFS", writefstest.TestFS(fsys.(writefs.WriteFS)))
	})

	t.Run("All methods returns factory error if any", func(t *testing.T) {

	})
}
