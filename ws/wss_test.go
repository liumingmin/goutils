package ws

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/liumingmin/goutils/utils/safego"

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
	const pullMsgFromDB = 1
	var createCometFunc = func(conn *Connection) Comet {
		comet := NewDefaultComet(conn, pullMsgFromDB, func(ctx context.Context, pullConn *Connection) {
			packet := GetPoolMessage(S2C_RESP)
			packet.PMsg().Data = []byte("first msg from db")
			pullConn.SendMsg(ctx, packet, nil)
		}, func(ctx context.Context, pullConn *Connection) {
			//msg from db...
			time.Sleep(time.Second * 1)

			packet := GetPoolMessage(S2C_RESP)
			packet.PMsg().Data = []byte("pull msg from db")
			pullConn.SendMsg(ctx, packet, nil)
		})
		return comet
	}

	e := gin.New()
	e.GET("/join", func(ctx *gin.Context) {
		connMeta := ConnectionMeta{
			UserId:   ctx.DefaultQuery("uid", ""),
			Typed:    0,
			DeviceId: "",
			Version:  0,
			Charset:  0,
		}
		_, err := AcceptGin(ctx, connMeta, DebugOption(true),
			SrvUpgraderCompressOption(true),
			CompressionLevelOption(2),
			ConnEstablishHandlerOption(func(conn *Connection) {
				log.Info(context.Background(), "server conn establish: %v", conn.Id())
				comet := createCometFunc(conn)
				safego.Go(comet.Pull)
			}),
			ConnClosingHandlerOption(func(conn *Connection) {
				log.Info(context.Background(), "server conn closing: %v", conn.Id())
			}),
			ConnClosedHandlerOption(func(conn *Connection) {
				log.Info(context.Background(), "server conn closed: %v", conn.Id())
			}),
			SrvPullChannelsOption([]int{pullMsgFromDB}))
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
		DebugOption(true),
		ClientIdOption("server1"),
		ClientDialWssOption(url, false),
		ClientDialCompressOption(true),
		CompressionLevelOption(2),
		ConnEstablishHandlerOption(func(conn *Connection) {
			log.Info(context.Background(), "client conn establish: %v", conn.Id())
		}),
		ConnClosingHandlerOption(func(conn *Connection) {
			log.Info(context.Background(), "client conn closing: %v", conn.Id())
		}),
		ConnClosedHandlerOption(func(conn *Connection) {
			log.Info(context.Background(), "client conn closed: %v", conn.Id())
		}),
	)
	log.Info(ctx, "%v", conn)
	time.Sleep(time.Second * 5)

	packet := GetPoolMessage(C2S_REQ)
	packet.PMsg().Data = []byte("client request")
	conn.SendMsg(context.Background(), packet, nil)

	time.Sleep(time.Second * 10)

	//client connect displace
	conn2, _ := DialConnect(context.Background(), url, http.Header{},
		DebugOption(true),
		ClientIdOption("server2"),
		ClientDialWssOption(url, false),
		ConnEstablishHandlerOption(func(conn *Connection) {
			log.Info(context.Background(), "client conn establish: %v", conn.Id())
		}),
	)
	time.Sleep(time.Second)
	packet = GetPoolMessage(C2S_REQ)
	packet.PMsg().Data = []byte("client request2")
	conn2.SendMsg(context.Background(), packet, nil)

	time.Sleep(time.Minute * 1)
}
