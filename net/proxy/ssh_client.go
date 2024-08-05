package proxy

import (
	"context"
	"fmt"
	"net"

	"github.com/liumingmin/goutils/log"

	"golang.org/x/crypto/ssh"
)

type SshClientConfig struct {
	config *ssh.ClientConfig
	server string
}

//address ssh地址 127.0.0.1:22
//user 连接用户 app
//password 密码

func NewPassSshClient(address, user, password string) (*SshClientConfig, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// use OpenSSH's known_hosts file if you care about host validation
			return nil
		},
	}

	client := &SshClientConfig{
		config: config,
		server: address,
	}

	return client, nil
}

//address ssh地址 127.0.0.1:22
//user 连接用户 app
//privateKeyBytes 私钥内容
//privateKeyPassword 私钥密码

func NewSshClient(address, user string, privateKeyBytes []byte, privateKeyPassword string) (*SshClientConfig, error) {
	signer, err := ssh.ParsePrivateKeyWithPassphrase(privateKeyBytes, []byte(privateKeyPassword))
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// use OpenSSH's known_hosts file if you care about host validation
			return nil
		},
	}

	client := &SshClientConfig{
		config: config,
		server: address,
	}

	return client, nil
}

func (s *SshClientConfig) SshDial() (*ssh.Client, error) {
	return ssh.Dial("tcp", s.server, s.config)
}

func (s *SshClientConfig) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	sshClient, err := s.SshDial()
	if err != nil {
		log.Error(ctx, "DialContext err: %v", err)
		return nil, fmt.Errorf("DialContext err: %v", err)
	}
	//todo sshClient close?
	return sshClient.Dial(network, address)
}
