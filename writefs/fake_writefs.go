package writefs

import (
	"io/fs"
	"testing/fstest"
)

type testWriteFS struct {
	fstest.MapFS
	expectedErr error
}

var _ WriteFS = testWriteFS{}

type testRemoveFS struct {
	testWriteFS
	removed string
}

var _ RemoveFS = &testRemoveFS{}

func (fsys *testRemoveFS) Remove(name string) error {
	fsys.removed = name
	return nil
}

type testMkDirFS struct {
	testWriteFS
	created string
}

var _ MkDirFS = &testMkDirFS{}

func (fsys *testMkDirFS) MkDir(name string, perm fs.FileMode) error {
	fsys.created = name
	return nil
}

type testFileWriter struct {
}

var _ FileWriter = testFileWriter{}

func (w testFileWriter) Close() error {
	return nil
}
func (w testFileWriter) Write(buf []byte) (int, error) {
	return len(buf), nil
}
func (w testFileWriter) Read(buf []byte) (int, error) {
	return len(buf), nil
}
func (w testFileWriter) Stat() (fs.FileInfo, error) {
	return nil, nil
}

// OpenFile ...
func (fsys testWriteFS) OpenFile(name string, flag int, perm fs.FileMode) (FileWriter, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{}
	}
	return testFileWriter{}, fsys.expectedErr
}
