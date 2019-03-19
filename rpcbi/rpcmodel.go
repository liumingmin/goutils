package rpcbi

const (
	PROTOCOL_HANDSHAKE     = 1
	PROTOCOL_HANDSHAKE_ACK = 2
)

const (
	PROTOCOL_VERSION = 1
)

type HandshakeReq struct {
	Id      string
	Version int
}

type HandshakeResp struct {
	Code int
	Msg  string
}
