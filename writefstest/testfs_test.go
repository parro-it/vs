package writefstest

import (
	"testing"

	"github.com/parro-it/vs/memfs"
	"github.com/parro-it/vs/osfs"
)

func TestMemFS(t *testing.T) {
	fsys := memfs.New()
	t.Run("writefstest.TestFS", TestFS(fsys))
}

func TestOSFS(t *testing.T) {
	fsys := osfs.DirWriteFS("/var/fixtures")
	t.Run("writefstest.TestFS", TestFS(fsys))
}
