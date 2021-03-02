package sshfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/parro-it/sshconfig"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// ConnectClient returns a functioning instance of *SSHFS
// using the given ssh.Client as transport layer.
func ConnectClient(root string, sshClient *ssh.Client) (*SSHFS, error) {
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	fsys := SSHFS{
		client: client,
		root:   root,
	}
	return &fsys, nil
}

type hostCfg struct {
	*ssh.ClientConfig
	HostPort string
}

var cfg map[string]*hostCfg

func initConfig() error {

	if cfg == nil {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		hosts, err := sshconfig.Parse(os.DirFS(home), ".ssh/config")
		if err != nil {
			return err
		}
		cfg = map[string]*hostCfg{}

		for _, host := range hosts {
			hostCfg, err := hostToCfg(host)
			if err != nil {
				return err
			}
			if hostCfg != nil {
				cfg[host.Host[0]] = hostCfg
			}
		}
	}

	return nil
}

func hostToCfg(host *sshconfig.SSHHost) (*hostCfg, error) {
	if host.Host != nil && len(host.Host) > 0 && host.Host[0] == "*" {
		return nil, nil
	}
	if host.IdentityFile == "" {
		return nil, nil
	}

	identityFile := host.IdentityFile
	if strings.Contains(identityFile, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		identityFile = strings.ReplaceAll(identityFile, "~", home)
	}

	key, err := privateSSHKey(identityFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read ssh key %s: %w", host.IdentityFile, err)
	}

	hostCfg := &hostCfg{
		ClientConfig: &ssh.ClientConfig{
			User:            host.User,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 5,
			Auth:            []ssh.AuthMethod{key},
		},
		HostPort: fmt.Sprintf("%s:%d", host.HostName, host.Port),
	}
	return hostCfg, nil
}

// ConnectFromConfig returns a functioning instance of *SSHFS
// using the info in ~/.ssh/config to create and connect an
// SSH transport layer.
// the Disconnect method of the SSHFS instance will disconnect the
// SSH connection too.
func ConnectFromConfig(root string,sshHostName string) (*SSHFS, error) {
	err := initConfig()
	if err != nil {
		return nil, err
	}
	return connect(root,cfg[sshHostName])
}

// Connect returns a functioning instance of *SSHFS
// using the given ssh.ClientConfig configuration to create
// and connect an SSH transport layer.
// the Disconnect method of the SSHFS instance will disconnect the
// SSH connection too.
func Connect(root string, config *sshconfig.SSHHost) (*SSHFS, error) {
	hostCfg, err := hostToCfg(config)
	if err != nil {
		return nil, err
	}
	if hostCfg == nil {
		return nil, fmt.Errorf("unvalid config provided")
	}
	return connect(root, hostCfg)
}

func connect(root string, config *hostCfg) (*SSHFS, error) {
	sshClient, err := ssh.Dial("tcp", config.HostPort, config.ClientConfig)
	if err != nil {
		return nil, err
	}

	client, err := ConnectClient(root, sshClient)
	client.ownedSSHCient = sshClient
	return client, err
}

// Disconnect ...
func (fsys *SSHFS) Disconnect() {
	fsys.client.Close()

	if fsys.ownedSSHCient != nil {
		fsys.ownedSSHCient.Close()
		fsys.ownedSSHCient = nil
	}
}

func privateSSHKey(path string) (ssh.AuthMethod, error) {
	privateKey, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(privateKey)

	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}
