package ws

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liumingmin/goutils/log"
)

func TestWssRun(t *testing.T) {
	InitServer() //server invoke 服务端调用
	InitClient() //client invoke 客户端调用
	ctx := context.Background()

	const (
		C2S_REQ  = 1
		S2C_RESP = 2
	)

	//server reg handler
	RegisterHandler(C2S_REQ, func(ctx context.Context, connection *Connection, message *P_MESSAGE) error {
		log.Info(ctx, "server recv: %v, %v", message.ProtocolId, string(message.Data))
		packet := GetPMessage()
		packet.ProtocolId = S2C_RESP
		packet.Data = []byte("server response")
		connection.SendMsg(ctx, packet, nil)
		return nil
	})

	//server start
	e := gin.New()
	e.GET("/join", func(ctx *gin.Context) {
		connMeta := ConnectionMeta{
			UserId:   ctx.DefaultQuery("uid", ""),
			Typed:    0,
			DeviceId: "",
			Version:  0,
			Charset:  0,
		}
		_, err := AcceptGin(ctx, connMeta, ConnectCbOption(&ConnectCb{connMeta.UserId}))
		if err != nil {
			log.Error(ctx, "Accept client connection failed. error: %v", err)
			return
		}
	})
	go e.Run(":8003")

	//client reg handler
	RegisterHandler(S2C_RESP, func(ctx context.Context, connection *Connection, message *P_MESSAGE) error {
		log.Info(ctx, "client recv: %v, %v", message.ProtocolId, string(message.Data))
		return nil
	})
	//client connect
	uid := "100"
	conn, _ := Connect(context.Background(), "server1", "ws://127.0.0.1:8003/join?uid="+uid, false, http.Header{})
	log.Info(ctx, "%v", conn)
	time.Sleep(time.Second * 5)

	packet := GetPMessage()
	packet.ProtocolId = C2S_REQ
	packet.Data = []byte("client request")
	conn.SendMsg(context.Background(), packet, nil)

	time.Sleep(time.Minute * 1)
}

type ConnectCb struct {
	Uid string
}

func (c *ConnectCb) ConnFinished(clientId string) {
	log.Debug(context.Background(), "%v connected", c.Uid)
}
func (c *ConnectCb) DisconnFinished(clientId string) {
	log.Debug(context.Background(), "%v disconnected", c.Uid)
}
