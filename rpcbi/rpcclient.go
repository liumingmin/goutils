package rpcbi

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"encoding/json"

	"errors"

	"github.com/hashicorp/yamux"
	"github.com/liumingmin/goutils/safego"
	"github.com/liumingmin/goutils/utils"
)

type RpcClient struct {
	rpcServer *rpc.Server
	*rpc.Client
	protocolFormat int

	ClientId string
}

func NewRpcClient(protocolFormat int, id string) *RpcClient {
	return &RpcClient{
		rpcServer:      rpc.NewServer(),
		ClientId:       id,
		protocolFormat: protocolFormat,
	}
}

func (c *RpcClient) doServer(sess *yamux.Session) {
	clientConn, err := sess.Accept()
	if err != nil {
		log.Panic(err)
		return
	}

	if c.protocolFormat == PROTOCOL_FORMAT_GOB {
		c.rpcServer.ServeConn(clientConn)
	} else if c.protocolFormat == PROTOCOL_FORMAT_JSON {
		c.rpcServer.ServeCodec(jsonrpc.NewServerCodec(clientConn))
	}
}

func (c *RpcClient) doClient(sess *yamux.Session) {
	clientConn, err := sess.Open()
	if err != nil {
		log.Panic(err)
		return
	}

	if c.protocolFormat == PROTOCOL_FORMAT_GOB {
		c.Client = rpc.NewClient(clientConn)
	} else if c.protocolFormat == PROTOCOL_FORMAT_JSON {
		c.Client = rpc.NewClientWithCodec(jsonrpc.NewClientCodec(clientConn))
	}
}

func (c *RpcClient) Start(conn net.Conn) error {
	err := c.handshake(conn)
	if err != nil {
		conn.Close()
		return err
	}

	sess, err := yamux.Client(conn, nil)
	if err != nil {
		log.Panic(err)
		return err
	}

	safego.Go(func() {
		c.doServer(sess)
	})
	c.doClient(sess)
	return nil
}

func (c *RpcClient) handshake(conn net.Conn) error {
	bs, _ := json.Marshal(&HandshakeReq{Version: PROTOCOL_VERSION, Id: c.ClientId})
	handshake := &utils.DataPacket{PacketHeader: utils.PacketHeader{ProtocolId: PROTOCOL_HANDSHAKE}, Data: bs}
	err := handshake.Pack(conn)
	if err != nil {
		return err
	}

	handshakeAck := &utils.DataPacket{}
	err = handshakeAck.Unpack(conn)
	if err != nil {
		return err
	}

	if handshakeAck.ProtocolId != PROTOCOL_HANDSHAKE_ACK {
		return errors.New("protocol not match")
	}

	resp := &HandshakeResp{}
	err = json.Unmarshal(handshakeAck.Data, resp)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return errors.New(resp.Msg)
	}

	return nil
}

func (c *RpcClient) RegisterService(name string, service interface{}) error {
	return c.rpcServer.RegisterName(name, service)
}

func (c *RpcClient) Close() {
	c.Client.Close()
}
