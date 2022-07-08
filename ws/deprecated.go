package ws

//deprecated
type IConnCallback interface {
	ConnFinished(clientId string)
	DisconnFinished(clientId string)
}

//deprecated
type IHeartbeatCallback interface {
	RecvPing(clientId string)
	RecvPong(clientId string) error
}

//deprecated
func ConnectCbOption(connCallback IConnCallback) ConnOption {
	return func(conn *Connection) {
		conn.connEstablishHandler = func(c *Connection) {
			connCallback.ConnFinished(c.Id())
		}

		conn.connClosingHandler = func(c *Connection) {
			connCallback.DisconnFinished(c.Id())
		}
	}
}

//deprecated
func HeartbeatCbOption(heartbeatCallback IHeartbeatCallback) ConnOption {
	return func(conn *Connection) {
		conn.recvPingHandler = func(c *Connection) {
			heartbeatCallback.RecvPing(c.Id())
		}

		conn.recvPongHandler = func(c *Connection) {
			heartbeatCallback.RecvPong(c.Id())
		}
	}
}
