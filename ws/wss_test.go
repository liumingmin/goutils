package ws

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/liumingmin/goutils/utils/safego"

	"github.com/gin-gonic/gin"
	"github.com/liumingmin/goutils/log"
)

func TestMessage(t *testing.T) {
	InitClient()

	//poolMsg := GetPoolMessage(int32(P_S2C_s2c_err_displace))
	//displace := poolMsg.DataMsg().(*P_DISPLACE)
	//displace.Ts = time.Now().Unix()
	//displace.OldIp = []byte("1")
	//displace.NewIp = []byte("2")
	//bs, _ := poolMsg.Marshal()
	//for _, b := range bs {
	//	fmt.Print(fmt.Sprintf("%v,", b))
	//}

	m := NewMessage().(*Message)
	err := m.Unmarshal([]byte{8, 255, 255, 255, 255, 255, 255, 255, 255, 255, 1, 18, 12, 10, 1, 49, 18, 1, 50, 24, 143, 186, 156, 152, 6})
	t.Log(err)
	t.Log(m.PMsg())
	t.Log(m.DataMsg())
}

func TestWssRun(t *testing.T) {
	//InitServer()
	InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke 服务端调用
	InitClient()                                                    //client invoke 客户端调用
	ctx := context.Background()

	const (
		C2S_REQ  = 1
		S2C_RESP = 2
	)
	const pullMsgFromDB = 1

	//server reg handler
	RegisterHandler(C2S_REQ, func(ctx context.Context, connection IConnection, message IMessage) error {
		log.Info(ctx, "server recv: %v, %v", message.PMsg().ProtocolId, string(message.PMsg().Data))
		packet := GetPoolMessage(S2C_RESP)
		packet.PMsg().Data = []byte("server response")
		connection.SendMsg(ctx, packet, nil)

		connection.SendPullNotify(ctx, pullMsgFromDB)
		return nil
	})

	rawMsg := &P_MESSAGE{}
	rawMsg.ProtocolId = S2C_RESP
	rawMsg.Data = []byte("common msg")

	commonMsg := NewMessage()
	commonMsg.PMsg().Data, _ = proto.Marshal(rawMsg)

	//server start
	var createSrvPullerFunc = func(conn IConnection, pullChannelId int) Puller {
		return NewDefaultPuller(conn, pullChannelId, func(ctx context.Context, pullConn IConnection) {
			packet := GetPoolMessage(S2C_RESP)
			packet.PMsg().Data = []byte("first msg from db")
			pullConn.SendMsg(ctx, packet, nil)
		}, func(ctx context.Context, pullConn IConnection) {
			//msg from db...
			time.Sleep(time.Second * 1)

			packet := GetPoolMessage(S2C_RESP)
			packet.PMsg().Data = []byte("pull msg from db")
			pullConn.SendMsg(ctx, packet, nil)

			pullConn.SendMsg(ctx, commonMsg, nil)
		})
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
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn establish: %v, %p", conn.Id(), conn)
				//在集群环境下，需要检查connId是否已经连接集群，如有需踢掉在集群其他节点建立的连接，可通过redis pub sub，其他节点收到通知调用KickClient
				//lastConnNodeId, lastConnMTs := GetClientTs(ctx, conn.Id())
				//if lastConnNodeId != "" && lastConnNodeId != config.NodeId && lastConnMTs < util.UtcMTs() {
				//	MqPublish(ctx, conn.Id(), conn.ClientIp())   //other node: ClientConnHub.Find(connId).DisplaceClientByIp(ctx, newIp)
				//}
				//RegisterConn() // save to redis
				puller := createSrvPullerFunc(conn, pullMsgFromDB)
				safego.Go(func() {
					puller.PullSend()
				})
			}),
			ConnClosingHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn closing: %v, %p", conn.Id(), conn)
			}),
			ConnClosedHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn closed: %v, %p", conn.Id(), conn)
			}),
			SrvPullChannelsOption([]int{pullMsgFromDB}))
		if err != nil {
			log.Error(ctx, "Accept client connection failed. error: %v", err)
			return
		}
	})
	go e.Run(":8003")

	//client reg handler
	RegisterHandler(S2C_RESP, func(ctx context.Context, connection IConnection, message IMessage) error {
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
		ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
			log.Info(ctx, "client conn establish: %v, %p", conn.Id(), conn)
		}),
		ConnClosingHandlerOption(func(ctx context.Context, conn IConnection) {
			log.Info(ctx, "client conn closing: %v, %p", conn.Id(), conn)
		}),
		ConnClosedHandlerOption(func(ctx context.Context, conn IConnection) {
			log.Info(ctx, "client conn closed: %v, %p", conn.Id(), conn)
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
		ClientDialCompressOption(true),
		CompressionLevelOption(2),
		ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
			log.Info(ctx, "client conn establish: %v, %p", conn.Id(), conn)
		}),
	)
	time.Sleep(time.Second)
	packet = GetPoolMessage(C2S_REQ)
	packet.PMsg().Data = []byte("client request2")
	conn2.SendMsg(context.Background(), packet, nil)

	clientConn, _ := ClientConnHub.Find("100-0-")
	clientConn.DisplaceClientByIp(ctx, "newip")
	for i := 0; i < 100; i++ {
		url := "ws://127.0.0.1:8003/join?uid=a" + strconv.Itoa(i)
		DialConnect(context.Background(), url, http.Header{},
			DebugOption(true),
			ClientIdOption(strconv.Itoa(i)),
			ClientDialWssOption(url, false),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				safego.Go(func() {
					time.Sleep(time.Second * 3)
					conn.KickServer()
				})
				log.Info(ctx, "client conn establish: %v", conn.Id())
			}),
		)
	}
	fmt.Println(len(ClientConnHub.ConnectionIds()))
	time.Sleep(time.Second * 5)
	for i := 0; i < 100; i++ {
		url := "ws://127.0.0.1:8003/join?uid=b" + strconv.Itoa(i)
		DialConnect(context.Background(), url, http.Header{},
			DebugOption(true),
			ClientIdOption("b"+strconv.Itoa(i)),
			ClientDialWssOption(url, false),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "client conn establish: %v", conn.Id())
			}),
		)
	}

	time.Sleep(time.Second * 5)

	fmt.Println(len(ClientConnHub.ConnectionIds()))

	time.Sleep(time.Minute * 1)
}
