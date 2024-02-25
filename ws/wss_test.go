package ws

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/liumingmin/goutils/utils/safego"

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

	//m := NewMessage().(*Message)
	//err := m.unmarshal([]byte{8, 255, 255, 255, 255, 255, 255, 255, 255, 255, 1, 18, 12, 10, 1, 49, 18, 1, 50, 24, 143, 186, 156, 152, 6})
	//t.Log(err)
	//t.Log(m.protocolId())
	//t.Log(m.DataMsg())
}

func TestWssRequestResponse(t *testing.T) {
	InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke
	InitClient()                                                    //client invoke
	ctx := context.Background()

	const (
		C2S_REQ  = 2
		S2C_RESP = 3
	)

	//server reg handler
	RegisterHandler(C2S_REQ, func(ctx context.Context, connection IConnection, message IMessage) error {
		log.Info(ctx, "server recv: %v, %v", message.GetSn(), string(message.GetData()))
		packet := GetPoolMessage(S2C_RESP)
		packet.SetData([]byte("server rpc resp info"))
		return connection.SendResponseMsg(ctx, packet, message.GetSn(), nil)
	})

	//server start
	// e := gin.New()
	// e.Static("/js", "js")
	// e.Static("/ts", "ts")

	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		connMeta := ConnectionMeta{
			UserId:   r.URL.Query().Get("uid"),
			Typed:    0,
			DeviceId: "",
			Version:  0,
			Charset:  0,
		}
		_, err := Accept(ctx, w, r, connMeta, DebugOption(true),
			SrvUpgraderCompressOption(true),
			CompressionLevelOption(2),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn establish: %v, %p", conn.Id(), conn)
			}),
			ConnClosingHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn closing: %v, %p", conn.Id(), conn)
			}),
			ConnClosedHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn closed: %v, %p", conn.Id(), conn)
			}),
		)
		if err != nil {
			log.Error(ctx, "Accept client connection failed. error: %v", err)
			return
		}
	})

	go http.ListenAndServe(":8003", nil)

	//client reg handler
	RegisterHandler(S2C_RESP, func(ctx context.Context, connection IConnection, message IMessage) error {
		log.Info(ctx, "client handler recv sn: %v, %v", message.GetSn(), string(message.GetData()))
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
	time.Sleep(time.Second * 1)

	for i := 0; i < 10; i++ {
		packet := GetPoolMessage(C2S_REQ)
		packet.SetData([]byte("client rpc req info"))
		resp, err := conn.SendRequestMsg(context.Background(), packet, nil)
		if err == nil {
			log.Info(ctx, "client recv: sn: %v, data: %v", resp.GetSn(), string(resp.GetData()))
		}
	}

	time.Sleep(time.Second * 5)
}

func TestWssRequestResponseWithTimeout(t *testing.T) {
	InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke
	InitClient()                                                    //client invoke
	ctx := context.Background()

	const (
		C2S_REQ_TIMEOUT  = 4
		S2C_RESP_TIMEOUT = 5
	)

	//server reg handler
	RegisterHandler(C2S_REQ_TIMEOUT, func(ctx context.Context, connection IConnection, message IMessage) error {
		log.Info(ctx, "server recv: %v, %v", message.GetSn(), string(message.GetData()))

		time.Sleep(time.Second * 4)
		packet := GetPoolMessage(S2C_RESP_TIMEOUT)
		packet.SetData([]byte("server rpc resp info timeout"))
		return connection.SendResponseMsg(ctx, packet, message.GetSn(), nil)
	})

	//server start
	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		connMeta := ConnectionMeta{
			UserId:   r.URL.Query().Get("uid"),
			Typed:    0,
			DeviceId: "",
			Version:  0,
			Charset:  0,
		}
		_, err := Accept(ctx, w, r, connMeta, DebugOption(true),
			SrvUpgraderCompressOption(true),
			CompressionLevelOption(2),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn establish: %v, %p", conn.Id(), conn)
			}),
			ConnClosingHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn closing: %v, %p", conn.Id(), conn)
			}),
			ConnClosedHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn closed: %v, %p", conn.Id(), conn)
			}),
		)
		if err != nil {
			log.Error(ctx, "Accept client connection failed. error: %v", err)
			return
		}
	})
	go http.ListenAndServe(":8003", nil)

	//client reg handler
	RegisterHandler(S2C_RESP_TIMEOUT, func(ctx context.Context, connection IConnection, message IMessage) error {
		log.Info(ctx, "client handler recv sn: %v, %v", message.GetSn(), string(message.GetData()))
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
	time.Sleep(time.Second * 1)

	toCtx, _ := context.WithTimeout(ctx, time.Second*2)
	packet := GetPoolMessage(C2S_REQ_TIMEOUT)
	packet.SetData([]byte("client rpc req info timeout"))
	resp, err := conn.SendRequestMsg(toCtx, packet, nil)
	if err == nil {
		log.Info(ctx, "client recv: sn: %v, data: %v", resp.GetSn(), string(resp.GetData()))
	} else {
		log.Error(ctx, "client recv err: %v", err)
	}

	time.Sleep(time.Second * 5)
}

func TestWssSendMessage(t *testing.T) {
	//InitServer()
	InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke
	InitClient()                                                    //client invoke
	ctx := context.Background()

	const (
		C2S_REQ  = 2
		S2C_RESP = 3
	)
	const pullMsgFromDB = 1

	//server reg handler
	RegisterHandler(C2S_REQ, func(ctx context.Context, connection IConnection, message IMessage) error {
		log.Info(ctx, "server recv: %v, %v", message.GetProtocolId(), string(message.GetData()))
		packet := GetPoolMessage(S2C_RESP)
		packet.SetData([]byte("server response"))
		connection.SendMsg(ctx, packet, nil)

		connection.SendPullNotify(ctx, pullMsgFromDB)
		return nil
	})

	commonMsg := NewMessage(S2C_RESP)
	commonMsg.SetData([]byte("common msg"))

	//server start
	var createSrvPullerFunc = func(conn IConnection, pullChannelId int) Puller {
		return NewDefaultPuller(conn, pullChannelId, func(ctx context.Context, pullConn IConnection) {
			packet := GetPoolMessage(S2C_RESP)
			packet.SetData([]byte("first msg from db"))
			pullConn.SendMsg(ctx, packet, nil)
		}, func(ctx context.Context, pullConn IConnection) {
			//msg from db...
			time.Sleep(time.Second * 1)

			packet := GetPoolMessage(S2C_RESP)
			packet.SetData([]byte("pull msg from db"))
			pullConn.SendMsg(ctx, packet, nil)

			pullConn.SendMsg(ctx, commonMsg, nil)
		})
	}

	//e := gin.New()
	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		connMeta := ConnectionMeta{
			UserId:   r.URL.Query().Get("uid"),
			Typed:    0,
			DeviceId: "",
			Version:  0,
			Charset:  0,
		}
		_, err := Accept(ctx, w, r, connMeta, DebugOption(true),
			SrvUpgraderCompressOption(true),
			CompressionLevelOption(2),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn establish: %v, %p", conn.Id(), conn)
				// In a cluster environment, it is necessary to check whether connId has already connected to the cluster.
				// If it is necessary to kick off connections established on other nodes in the cluster,
				// it can be done through Redis pub sub, and other nodes will receive a notification to call KickClient

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
	go http.ListenAndServe(":8003", nil)

	//client reg handler
	RegisterHandler(S2C_RESP, func(ctx context.Context, connection IConnection, message IMessage) error {
		log.Info(ctx, "client recv: %v, %v", message.GetProtocolId(), string(message.GetData()))
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
	time.Sleep(time.Second * 1)

	//send msg
	packet := GetPoolMessage(C2S_REQ)
	packet.SetData([]byte("client request"))
	conn.SendMsg(context.Background(), packet, nil)

	time.Sleep(time.Second * 1)

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

	time.Sleep(time.Second * 1)

	//send msg by conn2
	packet = GetPoolMessage(C2S_REQ)
	packet.SetData([]byte("client request2"))
	conn2.SendMsg(context.Background(), packet, nil)

	time.Sleep(time.Second * 2)
}

func TestWssDialConnect(t *testing.T) {
	InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke
	InitClient()                                                    //client invoke
	ctx := context.Background()

	//server start
	//e := gin.New()
	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		connMeta := ConnectionMeta{
			UserId:   r.URL.Query().Get("uid"),
			Typed:    0,
			DeviceId: "",
			Version:  0,
			Charset:  0,
		}
		_, err := Accept(ctx, w, r, connMeta, DebugOption(true),
			SrvUpgraderCompressOption(true),
			CompressionLevelOption(2),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn establish: %v, %p", conn.Id(), conn)
				// In a cluster environment, it is necessary to check whether connId has already connected to the cluster.
				// If it is necessary to kick off connections established on other nodes in the cluster,
				// it can be done through Redis pub sub, and other nodes will receive a notification to call KickClient

				//lastConnNodeId, lastConnMTs := GetClientTs(ctx, conn.Id())
				//if lastConnNodeId != "" && lastConnNodeId != config.NodeId && lastConnMTs < util.UtcMTs() {
				//	MqPublish(ctx, conn.Id(), conn.ClientIp())   //other node: ClientConnHub.Find(connId).DisplaceClientByIp(ctx, newIp)
				//}
				//RegisterConn() // save to redis
			}),
			ConnClosingHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn closing: %v, %p", conn.Id(), conn)
			}),
			ConnClosedHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn closed: %v, %p", conn.Id(), conn)
			}))
		if err != nil {
			log.Error(ctx, "Accept client connection failed. error: %v", err)
			return
		}
	})
	go http.ListenAndServe(":8003", nil)

	//client connect
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
				//log.Info(ctx, "client conn establish: %v", conn.Id())
			}),
		)
	}

	fmt.Println(len(ClientConnHub.ConnectionIds()))
	time.Sleep(time.Second * 3)

	//client connect again
	for i := 0; i < 100; i++ {
		url := "ws://127.0.0.1:8003/join?uid=a" + strconv.Itoa(i)
		DialConnect(context.Background(), url, http.Header{},
			DebugOption(true),
			ClientIdOption("b"+strconv.Itoa(i)),
			ClientDialWssOption(url, false),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				//log.Info(ctx, "client conn establish: %v", conn.Id())
			}),
		)
	}

	time.Sleep(time.Second * 3)

	fmt.Println(len(ClientConnHub.ConnectionIds()))
}

func TestMain(m *testing.M) {
	//ignore auto test
	if runtime.GOOS != "windows" {
		return
	}
	m.Run()
}
