package rpcbi

import (
	"encoding/json"
	"errors"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"

	"time"

	"github.com/hashicorp/yamux"
	"github.com/liumingmin/goutils/log4go"
	"github.com/liumingmin/goutils/safego"
	"github.com/liumingmin/goutils/utils"
)

type RpcSession struct {
	socketConn net.Conn
	*rpc.Client
}

func (s *RpcSession) Close() {
	s.Client.Close()
	s.socketConn.Close()
}

type RpcServer struct {
	sessions       sync.Map // map[id]*Session
	rpcServer      *rpc.Server
	connCallback   ConnCallback
	protocolFormat int
}

func NewRpcServer(protocolFormat int, params ...interface{}) *RpcServer {
	server := &RpcServer{
		rpcServer:      rpc.NewServer(),
		protocolFormat: protocolFormat,
	}

	if len(params) > 0 {
		server.connCallback = params[0].(ConnCallback)
	}

	return server
}

func (s *RpcServer) newRpcSession(socketConn, sessionConn net.Conn) *RpcSession {
	session := &RpcSession{
		socketConn: socketConn,
	}

	if s.protocolFormat == PROTOCOL_FORMAT_GOB {
		session.Client = rpc.NewClient(sessionConn)
	} else if s.protocolFormat == PROTOCOL_FORMAT_JSON {
		session.Client = rpc.NewClientWithCodec(jsonrpc.NewClientCodec(sessionConn))
	}
	return session
}

func (s *RpcServer) serveRpc(sess *yamux.Session) {
	conn, err := sess.Accept()
	if err != nil {
		log4go.Error("Session accept connection failed.error: %v", err)
		return
	}

	if s.protocolFormat == PROTOCOL_FORMAT_GOB {
		s.rpcServer.ServeConn(conn)
	} else if s.protocolFormat == PROTOCOL_FORMAT_JSON {
		s.rpcServer.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func (c *RpcServer) doHandshake(conn net.Conn) (*HandshakeReq, error) {
	req := &utils.DataPacket{}
	err := req.Unpack(conn)
	if err != nil {
		return nil, err
	}

	if req.ProtocolId != PROTOCOL_HANDSHAKE {
		return nil, errors.New("protocol not handshake")
	}

	handReq := &HandshakeReq{}
	err = json.Unmarshal(req.Data, handReq)
	if err != nil {
		return nil, err
	}

	if handReq.Version != PROTOCOL_VERSION {
		return nil, errors.New("version not match")
	}

	return handReq, nil
}

func (c *RpcServer) handshake(conn net.Conn) (*HandshakeReq, error) {
	req, err := c.doHandshake(conn)
	var resp *HandshakeResp
	if err == nil {
		resp = &HandshakeResp{Code: 0, Msg: "ok"}
	} else {
		resp = &HandshakeResp{Code: -1, Msg: err.Error()}
	}

	bs, _ := json.Marshal(resp)
	handshake := &utils.DataPacket{PacketHeader: utils.PacketHeader{ProtocolId: PROTOCOL_HANDSHAKE_ACK}, Data: bs}
	err1 := handshake.Pack(conn)
	if err1 != nil {
		err = err1
	}
	return req, err
}

func (s *RpcServer) handleConn(conn net.Conn) {
	var req *HandshakeReq
	var err error
	ok := utils.AsyncInvokeWithTimeout(time.Second*5, func() {
		req, err = s.handshake(conn)
	})

	if !ok {
		log4go.Info("client handshake timeout:%s\n", conn.RemoteAddr())
		conn.Close()
		return
	}

	if err != nil {
		conn.Close()
		return
	}

	defer conn.Close()

	id := req.Id
	log4go.Info("client connected:%s,%s\n", id, conn.RemoteAddr())

	sess, err := yamux.Server(conn, nil) //default config  keepalive ...
	if err != nil {
		log4go.Error("Create session failed.error: %v", err)
		return
	}

	clientConn, err := sess.Open()
	if err != nil {
		log4go.Error("Open session failed.error: %v", err)
		return
	}
	session := s.newRpcSession(conn, clientConn)
	s.sessions.Store(id, session)
	if s.connCallback != nil {
		s.connCallback.ConnFinished(id)
	}

	s.serveRpc(sess)

	s.sessions.Delete(id)

	if s.connCallback != nil {
		s.connCallback.DisconnFinished(id)
	}

	log4go.Info("client disconnected:%s,%s\n", id, conn.RemoteAddr())
}

func (s *RpcServer) RegisterService(name string, service interface{}) error {
	return s.rpcServer.RegisterName(name, service)
}

func (s *RpcServer) RangeSession(f func(id string, sess *RpcSession)) {
	s.sessions.Range(func(k, v interface{}) bool {
		f(k.(string), v.(*RpcSession))
		return true
	})
}

func (s *RpcServer) GetSession(id string) *RpcSession {
	if v, ok := s.sessions.Load(id); ok {
		return v.(*RpcSession)
	}
	return nil
}

func (s *RpcServer) Serve(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log4go.Error("Accept connection failed.error: %v", err)
			continue
		}
		safego.Go(func() {
			s.handleConn(conn)
		})
	}
}
