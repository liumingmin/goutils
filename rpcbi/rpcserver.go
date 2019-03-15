package rpcbi

import (
	"net"
	"net/rpc"

	"sync"

	"unsafe"

	"github.com/hashicorp/yamux"
	"github.com/liumingmin/goutils/log4go"
	"github.com/liumingmin/goutils/safego"
)

type RpcServer struct {
	*Server
	sessions    map[unsafe.Pointer]string
	rpcSessions map[string]*yamux.Session
	rpcClients  map[string]*rpc.Client
	lock        sync.Mutex
	version     int
}

func NewRpcServer(version int) (*RpcServer, error) {
	s := NewServer()
	rpcServer := &RpcServer{
		Server:      s,
		version:     version,
		sessions:    make(map[unsafe.Pointer]string),
		rpcSessions: make(map[string]*yamux.Session),
		rpcClients:  make(map[string]*rpc.Client),
	}
	s.RegisterName("server", &ServerHandshake{rpcServer})
	return rpcServer, nil
}

func (s *RpcServer) Start(network, address string) {
	lis, _ := net.Listen(network, address)
	for {
		conn, err := lis.Accept()
		if err != nil {
			log4go.Error("rpc.Serve: accept:", err.Error())
			continue
		}

		session, err := yamux.Server(conn, nil)
		if err != nil {
			log4go.Error("yamux.Server:", err.Error())
			continue
		}

		stream, err := session.Accept()
		if err != nil {
			log4go.Error("session.Accept:", err.Error())
			continue
		}

		s.lock.Lock()
		s.sessions[unsafe.Pointer(session)] = ""
		s.lock.Unlock()

		log4go.Debug("connected!")
		safego.Go(func() {
			s.ServeConn(stream)
			log4go.Debug("run done")

			session.Close()

			s.lock.Lock()
			defer s.lock.Unlock()
			sp := unsafe.Pointer(session)
			clientId := s.sessions[sp]
			delete(s.sessions, sp)

			if clientId != "" {
				delete(s.rpcSessions, clientId)
				delete(s.rpcClients, clientId)
			}

			log4go.Debug("disconnected! %s", clientId)
		})
	}
}

type ServerHandshake struct {
	server *RpcServer
}

func (s *ServerHandshake) Handshake(args *HandshakeReq, reply *HandshakeResp) error {
	return nil
}

func (s *ServerHandshake) DoHandshake(args *HandshakeReq, reply *HandshakeResp, conn net.Conn) error {
	log4go.Debug("start handshake %v", args)

	if args.Version != s.server.version {
		*reply = HandshakeResp{Code: -1, Msg: "version not match"}
		return nil
	}

	return nil

	s.server.lock.Lock()
	defer s.server.lock.Unlock()

	if old, ok := s.server.rpcClients[args.Id]; ok {
		old.Close()
		delete(s.server.rpcClients, args.Id)
	}

	if oldSession, ok2 := s.server.rpcSessions[args.Id]; ok2 {
		oldSession.Close()
		delete(s.server.rpcSessions, args.Id)
	}

	session := conn.(*yamux.Stream).Session()
	stream, err := session.Open()
	if err != nil {
		log4go.Error("handshake error: %v", err)
		return err
	}

	rpcClient := rpc.NewClient(stream)

	s.server.rpcSessions[args.Id] = session
	s.server.rpcClients[args.Id] = rpcClient
	s.server.sessions[unsafe.Pointer(session)] = args.Id

	*reply = HandshakeResp{Code: 0, Msg: "handshake ok"}

	log4go.Debug("handshake ok: %v", err)
	return nil
}
