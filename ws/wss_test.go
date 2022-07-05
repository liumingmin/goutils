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
	RegisterHandler(C2S_REQ, func(ctx context.Context, connection *Connection, message *Message) error {
		log.Info(ctx, "server recv: %v, %v", message.PMsg().ProtocolId, string(message.PMsg().Data))
		packet := GetPoolMessage(S2C_RESP)
		packet.PMsg().Data = []byte("server response")
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
		_, err := AcceptGin(ctx, connMeta, ConnectCbOption(&ConnectCb{connMeta.UserId}),
			SrvUpgraderCompressOption(true),
			CompressionLevelOption(1))
		if err != nil {
			log.Error(ctx, "Accept client connection failed. error: %v", err)
			return
		}
	})
	go e.Run(":8003")

	//client reg handler
	RegisterHandler(S2C_RESP, func(ctx context.Context, connection *Connection, message *Message) error {
		log.Info(ctx, "client recv: %v, %v", message.PMsg().ProtocolId, string(message.PMsg().Data))
		return nil
	})
	//client connect
	uid := "100"
	url := "ws://127.0.0.1:8003/join?uid=" + uid
	conn, _ := DialConnect(context.Background(), url, http.Header{},
		ClientIdOption("server1"),
		ClientDialWssOption(url, false),
		ClientDialCompressOption(true),
		CompressionLevelOption(2),
	)
	log.Info(ctx, "%v", conn)
	time.Sleep(time.Second * 5)

	packet := GetPoolMessage(C2S_REQ)
	packet.PMsg().Data = []byte("client request")
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

func TestBenchmarkWssRun(t *testing.T) {
	InitServer() //server invoke 服务端调用
	InitClient() //client invoke 客户端调用
	const (
		C2S_REQ  = 1
		S2C_RESP = 2
	)

	var reqBytes = []byte("client request")
	var respBytes = []byte("server response")

	RegisterDataMsgType(C2S_REQ, &P_MESSAGE{})
	RegisterDataMsgType(S2C_RESP, &P_MESSAGE{})
	ctx := context.Background()

	//server reg handler
	RegisterHandler(C2S_REQ, func(ctx context.Context, connection *Connection, message *Message) error {
		//log.Info(ctx, "server recv: %v, %v", message.PMsg().ProtocolId, string(message.PMsg().Data))
		packet := GetPoolMessage(S2C_RESP)
		dataMsg := packet.DataMsg().(*P_MESSAGE)
		dataMsg.Data = respBytes
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
	RegisterHandler(S2C_RESP, func(ctx context.Context, connection *Connection, message *Message) error {
		//log.Info(ctx, "client recv: %v, %v", message.PMsg().ProtocolId, string(message.PMsg().Data))
		return nil
	})
	//client connect
	uid := "100"
	conn, _ := Connect(context.Background(), "server1", "ws://127.0.0.1:8003/join?uid="+uid, false, http.Header{})
	log.Info(ctx, "%v", conn)
	time.Sleep(time.Second * 5)

	for i := 0; i < 100000; i++ {
		packet := GetPoolMessage(C2S_REQ)
		dataMsg := packet.DataMsg().(*P_MESSAGE)
		dataMsg.Data = reqBytes
		conn.SendMsg(context.Background(), packet, nil)
		time.Sleep(time.Millisecond * 50)
	}

	time.Sleep(time.Minute * 3)
}
