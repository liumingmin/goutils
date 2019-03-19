package rpcbi

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"

	"encoding/json"

	"errors"

	"github.com/hashicorp/yamux"
	"github.com/liumingmin/goutils/safego"
	"github.com/liumingmin/goutils/utils"
)

type RpcJsonServer struct {
	clientid  int32
	sessions  sync.Map // map[id]*Session
	rpcServer *rpc.Server
}

func NewRpcJsonServer() *RpcJsonServer {
	return &RpcJsonServer{
		rpcServer: rpc.NewServer(),
	}
}

type RpcJsonSession struct {
	masterConn net.Conn
	*rpc.Client
}

func NewRpcJsonSession(masterConn, clientConn net.Conn) *RpcJsonSession {
	return &RpcJsonSession{
		masterConn: masterConn,
		Client:     rpc.NewClientWithCodec(jsonrpc.NewClientCodec(clientConn)),
	}
}

func (s *RpcJsonSession) Close() {
	s.Client.Close()
	s.masterConn.Close()
}

func (s *RpcJsonSession) ConnFinished() {

}

func (s *RpcJsonSession) DisconnFinished() {

}

func (s *RpcJsonServer) serveRpc(sess *yamux.Session) {
	conn, err := sess.Accept()
	if err != nil {
		log.Print(err)
		return
	}
	s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
}

func (c *RpcJsonServer) doHandshake(conn net.Conn) (*HandshakeReq, error) {
	req := &utils.ControlPacket{}
	req.Unpack(conn)

	if req.ProtocolId != PROTOCOL_HANDSHAKE {
		return nil, errors.New("protocol not handshake")
	}

	handReq := &HandshakeReq{}
	err := json.Unmarshal(req.Data, handReq)
	if err != nil {
		return nil, err
	}

	if handReq.Version != PROTOCOL_VERSION {
		return nil, errors.New("version not match")
	}

	return handReq, nil
}

func (c *RpcJsonServer) handshake(conn net.Conn) (*HandshakeReq, error) {
	req, err := c.doHandshake(conn)
	var resp *HandshakeResp
	if err == nil {
		resp = &HandshakeResp{Code: 0, Msg: "ok"}
	} else {
		resp = &HandshakeResp{Code: -1, Msg: err.Error()}
	}

	bs, _ := json.Marshal(resp)
	handshake := &utils.ControlPacket{ProtocolId: PROTOCOL_HANDSHAKE_ACK, Data: bs}
	handshake.Pack(conn)

	return req, err
}

func (s *RpcJsonServer) handleConn(conn net.Conn) {
	req, err := s.handshake(conn)
	if err != nil {
		conn.Close()
		return
	}

	defer conn.Close()

	id := req.Id
	log.Printf("allocated %s for %s", id, conn.RemoteAddr())

	sess, err := yamux.Server(conn, nil)
	if err != nil {
		log.Print(err)
		return
	}

	clientConn, err := sess.Open()
	if err != nil {
		log.Print(err)
		return
	}
	session := NewRpcJsonSession(conn, clientConn)
	s.sessions.Store(id, session)
	session.ConnFinished()

	s.serveRpc(sess)

	s.sessions.Delete(id)
	session.DisconnFinished()
	log.Printf("%s(%s) closed connection", conn.RemoteAddr(), id)
}

func (s *RpcJsonServer) RegisterService(name string, service interface{}) error {
	return s.rpcServer.RegisterName(name, service)
}

func (s *RpcJsonServer) RangeSession(f func(id int32, sess *RpcJsonSession)) {
	s.sessions.Range(func(k, v interface{}) bool {
		f(k.(int32), v.(*RpcJsonSession))
		return true
	})
}

func (s *RpcJsonServer) Serve(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		safego.Go(func() {
			s.handleConn(conn)
		})
	}
}
