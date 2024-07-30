package proxy

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
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
	signer, err := signerFromPem(privateKeyBytes, []byte(privateKeyPassword))
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

func signerFromPem(pemBytes []byte, password []byte) (ssh.Signer, error) {
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, errors.New("pem decode failed, no key found")
	}

	// handle encrypted key
	if x509.IsEncryptedPEMBlock(pemBlock) {
		// decrypt PEM
		var err error
		pemBlock.Bytes, err = x509.DecryptPEMBlock(pemBlock, password)
		if err != nil {
			return nil, fmt.Errorf("decrypting PEM block failed %w", err)
		}

		// get RSA, EC or DSA key
		key, err := parsePemBlock(pemBlock)
		if err != nil {
			return nil, err
		}

		// generate signer instance from key
		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			return nil, fmt.Errorf("creating signer from encrypted key failed %w", err)
		}

		return signer, nil
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing plain private key failed %w", err)
	}

	return signer, nil
}

func parsePemBlock(block *pem.Block) (interface{}, error) {
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing PKCS private key failed %v", err)
		} else {
			return key, nil
		}
	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing EC private key failed %v", err)
		} else {
			return key, nil
		}
	case "DSA PRIVATE KEY":
		key, err := ssh.ParseDSAPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parsing DSA private key failed %v", err)
		} else {
			return key, nil
		}
	default:
		return nil, fmt.Errorf("parsing private key failed, unsupported key type %q", block.Type)
	}
}
