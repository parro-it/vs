package sshfs

import (
	"testing"

	"github.com/mikkeloscar/sshconfig"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestSSHFS(t *testing.T) {
	hostCfg := &sshconfig.SSHHost{
		IdentityFile: "~/.ssh/fake-host",
		User:         "andrea.parodi",
		Port:         2222,
		HostName:     "localhost",
	}
	/*
		t.Run("writefstest.TestFS", func(t *testing.T) {
			fsys, err := ConnectFromConfig("/var/fixtures", "fakehost")
			assert.NoError(t, err)
			t.Run("Pass writefstest.TestFS", writefstest.TestFS(fsys))
			fsys.Disconnect()
		})
	*/
	t.Run("Connection", func(t *testing.T) {
		t.Run("can be created from an ssh config hostname", func(t *testing.T) {
			fsys, err := ConnectFromConfig("/var/fixtures", "fakehost")
			assert.NoError(t, err)
			assert.NotNil(t, fsys)
			fsys.Disconnect()
		})

		t.Run("can be created from an sshconfig struct", func(t *testing.T) {
			fsys, err := Connect("/var/fixtures", hostCfg)
			assert.NoError(t, err)
			assert.NotNil(t, fsys)
			fsys.Disconnect()
		})

		t.Run("can be created from an ssh client", func(t *testing.T) {
			config, err := hostToCfg(hostCfg)
			assert.NoError(t, err)
			sshClient, err := ssh.Dial("tcp", config.HostPort, config.ClientConfig)
			assert.NoError(t, err)

			fsys, err := ConnectClient("/var/fixtures", sshClient)
			assert.NoError(t, err)
			assert.NotNil(t, fsys)

			fsys.Disconnect()
			sshClient.Close()
		})
	})
	/*
		t.Run("can open files for write", func(t *testing.T) {
			fsys, err := Connect("/var/fixtures", hostCfg)
			f, err := writefs.OpenFile(fsys, "prova", os.O_WRONLY|os.O_CREATE, fs.FileMode(0664))

			buf := []byte(time.Now().Format(time.RFC3339Nano))
			n, err := f.Write(buf)
			assert.NoError(t, err)
			assert.Equal(t, n, len(buf))
			f.Close()

			actual, err := fs.ReadFile(fsys, "prova")
			assert.NoError(t, err)
			assert.Equal(t, buf, actual)

		})

		t.Run("creates directories with OpenFile", func(t *testing.T) {
			fsys, err := Connect("/tmp", hostCfg)
			assert.NoError(t, err)
			fsys.client.RemoveDirectory("/tmp/adir")
			fsys.client.Remove("/tmp/adir")
			f, err := fsys.OpenFile("adir", os.O_CREATE, fs.FileMode(0755)|fs.ModeDir)
			assert.NoError(t, err)
			assert.Nil(t, f)
			info, err := fsys.client.Stat("/tmp/adir")
			assert.NoError(t, err)
			assert.True(t, info.IsDir())
			err = fsys.client.RemoveDirectory("/tmp/adir")
			assert.NoError(t, err)
		})

		t.Run("remove directories with OpenFile", func(t *testing.T) {
			fsys, err := Connect("/tmp", hostCfg)
			fsys.client.MkdirAll("/tmp/adir")
			info, err := fsys.client.Stat("/tmp/adir")
			assert.NoError(t, err)
			assert.True(t, info.IsDir())

			f, err := fsys.OpenFile("adir", os.O_TRUNC, 0)
			assert.NoError(t, err)
			assert.Nil(t, f)

			info, err = fsys.client.Stat("/tmp/adir")
			assert.Error(t, err)
			assert.True(t, os.IsNotExist(err))
			assert.Nil(t, info)

			fsys.client.RemoveDirectory("/tmp/adir")

		})
		t.Run("can open files", func(t *testing.T) {
			fsys, err := Connect("/var/fixtures", hostCfg)
			f, err := fsys.Open("ciao.txt")
			assert.NoError(t, err)
			assert.NotNil(t, f)
			buf := make([]byte, 100)
			n, err := f.Read(buf)
			assert.Equal(t, 5, n)
			assert.Equal(t, "ciao\n", string(buf[:n]))
			f.Close()

			fsys.Disconnect()
		})

		t.Run("can stat files", func(t *testing.T) {
			fsys, err := Connect("/var/fixtures", hostCfg)
			info, err := fs.Stat(fsys, "ciao.txt")
			assert.NoError(t, err)
			assert.NotNil(t, info)
			assert.Equal(t, "ciao.txt", info.Name())

			fsys.Disconnect()
		})

		t.Run("can read files", func(t *testing.T) {
			fsys, err := Connect("/var/fixtures", hostCfg)
			buf, err := fs.ReadFile(fsys, "ciao.txt")
			assert.NoError(t, err)
			assert.NotNil(t, buf)
			assert.Equal(t, "ciao\n", string(buf))

			fsys.Disconnect()
		})

		t.Run("can readdir", func(t *testing.T) {
			fsys, err := Connect("/var/fixtures", hostCfg)
			files, err := fs.ReadDir(fsys, "new-dir")
			assert.NoError(t, err)
			assert.Equal(t, 4, len(files))
			assert.Equal(t, "file1.txt", files[0].Name())

			fsys.Disconnect()
		})

		t.Run("pass TestFS", func(t *testing.T) {
			fsys, err := Connect("/var/fixtures", hostCfg)
			assert.NoError(t, err)
			err = fstest.TestFS(fsys, "new-dir", "new-dir/file1.txt")
			assert.NoError(t, err)
		})
	*/
}
