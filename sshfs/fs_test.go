package sshfs

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/parro-it/sshconfig"
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

}
