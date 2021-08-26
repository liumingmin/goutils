package proxy

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	*ssh.Client
	config *ssh.ClientConfig
	server string
}

//address ssh地址 127.0.0.1:22
//user 连接用户 app
//privateKeyBytes 私钥内容
//privateKeyPassword 私钥密码

func NewSshClient(address, user string, privateKeyBytes []byte, privateKeyPassword string) (*SshClient, error) {
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

	client := &SshClient{
		config: config,
		server: address,
	}

	return client, nil
}

func (s *SshClient) Connect() error {
	var err error
	s.Client, err = ssh.Dial("tcp", s.server, s.config)
	return err
}

func signerFromPem(pemBytes []byte, password []byte) (ssh.Signer, error) {
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		return nil, errors.New("Pem decode failed, no key found")
	}

	// handle encrypted key
	if x509.IsEncryptedPEMBlock(pemBlock) {
		// decrypt PEM
		var err error
		pemBlock.Bytes, err = x509.DecryptPEMBlock(pemBlock, password)
		if err != nil {
			return nil, fmt.Errorf("Decrypting PEM block failed %v", err)
		}

		// get RSA, EC or DSA key
		key, err := parsePemBlock(pemBlock)
		if err != nil {
			return nil, err
		}

		// generate signer instance from key
		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			return nil, fmt.Errorf("Creating signer from encrypted key failed %v", err)
		}

		return signer, nil
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("Parsing plain private key failed %v", err)
	}

	return signer, nil
}

func parsePemBlock(block *pem.Block) (interface{}, error) {
	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Parsing PKCS private key failed %v", err)
		} else {
			return key, nil
		}
	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Parsing EC private key failed %v", err)
		} else {
			return key, nil
		}
	case "DSA PRIVATE KEY":
		key, err := ssh.ParseDSAPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("Parsing DSA private key failed %v", err)
		} else {
			return key, nil
		}
	default:
		return nil, fmt.Errorf("Parsing private key failed, unsupported key type %q", block.Type)
	}
}
