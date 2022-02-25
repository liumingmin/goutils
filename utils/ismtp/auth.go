package ismtp

import (
	"bytes"
	"errors"
	"net/smtp"
)

type plainAuth struct {
	identity, username, password string
	host                         string
}

func PlainAuth(identity, username, password, host string) smtp.Auth {
	return &plainAuth{identity, username, password, host}
}

func (a *plainAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	//if !server.TLS {
	//	return "", nil, errors.New("unencrypted connection")
	//}
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	resp := []byte(a.identity + "\x00" + a.username + "\x00" + a.password)
	return "PLAIN", resp, nil
}

func (a *plainAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		// We've already sent everything.
		return nil, errors.New("unexpected server challenge")
	}
	return nil, nil
}

type loginAuth struct {
	username, password string
	host               string
}

func LoginAuth(username, password, host string) smtp.Auth {
	return &loginAuth{username, password, host}
}

func (a loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}
	return "LOGIN", nil, nil
}

func (a loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		if bytes.EqualFold([]byte("username:"), fromServer) {
			return []byte(a.username), nil
		} else if bytes.EqualFold([]byte("password:"), fromServer) {
			return []byte(a.password), nil
		}
	}
	return nil, nil
}
