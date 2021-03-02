module github.com/parro-it/vs

go 1.16

require (
	github.com/parro-it/sshconfig v0.0.0-20200912130257-f8464e2038cf
	github.com/pkg/sftp v1.12.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
)

replace github.com/parro-it/sshconfig => ../sshconfig
