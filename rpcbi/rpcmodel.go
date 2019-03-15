package rpcbi

type HandshakeReq struct {
	Id      string
	Version int
}

type HandshakeResp struct {
	Code int
	Msg  string
}
