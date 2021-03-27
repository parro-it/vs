package syncfs

import (
	"testing"

	"github.com/parro-it/vs/memfs"
	"github.com/parro-it/vs/writefs"
	"github.com/parro-it/vs/writefstest"
)

func TestSyncFS(t *testing.T) {
	t.Run("pass writefstest.TestFS", func(t *testing.T) {
		fsys := New(memfs.New())
		t.Run("Pass writefstest.TestFS", writefstest.TestFS(fsys.(writefs.WriteFS)))
	})

	t.Run("All methods returns factory error if any", func(t *testing.T) {

	})
}
