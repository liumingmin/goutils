package rpcbi

import (
	"net"
	"net/rpc"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/liumingmin/goutils/log4go"
	"github.com/liumingmin/goutils/safego"
)

type RpcClient struct {
	*rpc.Client
	*rpc.Server
	addr    string
	tcp     *net.TCPConn
	id      string
	version int
}

func (c *RpcClient) Handshake(id string, version int) (*HandshakeResp, error) {
	c.id = id
	c.version = version

	resp := &HandshakeResp{}
	err := c.Call("server.Handshake", &HandshakeReq{Id: id, Version: version},
		resp)

	return resp, err
}

func NewRpcClient(addr string, keepalive time.Duration) (*RpcClient, error) {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log4go.Error("connect addr error: %v", err)
		return nil, err
	}

	tcp := conn.(*net.TCPConn)
	//if err = tcp.SetKeepAlive(true); err != nil {
	//	return nil, err
	//}
	//
	//if err = tcp.SetKeepAlivePeriod(keepalive); err != nil {
	//	return nil, err
	//}
	//
	//if err = tcp.SetLinger(0); err != nil {
	//	return nil, err
	//}

	session, err := yamux.Client(conn, nil)
	if err != nil {
		log4go.Error("connect addr error: %v", err)
		return nil, err
	}

	// Open a new stream
	stream, err := session.Open()
	if err != nil {
		log4go.Error("connect addr error: %v", err)
		return nil, err
	}

	client := rpc.NewClient(stream)
	server := rpc.NewServer()

	safego.Go(func() {
		server.ServeConn(stream)
	})

	return &RpcClient{
		Client: client,
		Server: server,
		addr:   addr,
		tcp:    tcp,
	}, nil
}
