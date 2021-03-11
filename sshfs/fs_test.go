package sshfs

import (
	"testing"

	"github.com/mikkeloscar/sshconfig"
	"github.com/parro-it/vs/writefstest"
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

	t.Run("writefstest.TestFS", func(t *testing.T) {
		fsys, err := ConnectFromConfig("/var/fixtures", "fakehost")
		assert.NoError(t, err)

		sess, err := fsys.ownedSSHCient.NewSession()
		defer sess.Close()
		assert.NoError(t, err)
		err = sess.Run("rm -rf /var/fixtures/dir1 /var/fixtures/dirempty")
		assert.NoError(t, err)

		t.Run("Pass writefstest.TestFS", writefstest.TestFS(fsys))
		fsys.Disconnect()
	})

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

}
