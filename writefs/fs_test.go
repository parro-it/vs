package writefs

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func fixtureFile(name string) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot retrieve the source file path")
	}

	file = filepath.Dir(filepath.Dir(file))

	return path.Join(file, "fixtures", name)
}

func TestWriteFS(t *testing.T) {

	t.Run("defaults to Open", func(t *testing.T) {
		data := []byte{0xca, 0xfe, 0xba, 0xbe}

		memfs2 := fstest.MapFS{
			"adir2/afile2": &fstest.MapFile{Data: data},
		}

		f, err := OpenFile(memfs2, "adir2/afile2", os.O_RDONLY, fs.FileMode(0))
		defer f.Close()
		assert.NoError(t, err)
		buf := make([]byte, len(data))
		n, err := f.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, n, len(data))
		assert.Equal(t, data, buf)
	})

}
