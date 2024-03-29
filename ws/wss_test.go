package ws

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils/safego"
	"go.uber.org/zap/zapcore"
)

func TestWssRequestResponse(t *testing.T) {
	log.SetLogLevel(zapcore.WarnLevel) //if you need info log, commmet this line

	InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke
	InitClient()                                                    //client invoke
	ctx := context.Background()

	const (
		C2S_REQ  = 2
		S2C_RESP = 3
	)

	serverResp := "server rpc resp info"

	//server reg handler
	RegisterHandler(C2S_REQ, func(ctx context.Context, connection IConnection, message IMessage) error {
		log.Info(ctx, "server recv: %v, %v", message.GetSn(), string(message.GetData()))
		packet := GetPoolMessage(S2C_RESP)
		packet.SetData([]byte(serverResp))
		return connection.SendResponseMsg(ctx, packet, message.GetSn(), nil)
	})

	//server start
	handler := http.NewServeMux()
	handler.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
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
			NetMaxFailureRetryOption(int(time.Millisecond)*500),
			NetReadWaitOption(3*time.Second),
			NetWriteWaitOption(3*time.Second),
			NetTemporaryWaitOption(time.Millisecond*500),
			MaxMessageBytesSizeOption(1024*1024*32),
			SrvCheckOriginOption(func(r *http.Request) bool {
				return true
			}),
		)
		if err != nil {
			log.Error(ctx, "Accept client connection failed. error: %v", err)
			return
		}
	})

	go http.ListenAndServe(":8003", handler)

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
		ClientDialHandshakeTimeoutOption(time.Second*5),
		ClientDialRetryOption(2, time.Second*2),
		NetReadWaitOption(3*time.Second),
		NetWriteWaitOption(3*time.Second),
		ClientDialConnFailedHandlerOption(func(ctx context.Context, conn IConnection) {
			log.Error(ctx, "clien conn failed")
		}),
		ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
			log.Info(ctx, "client conn establish: %v, %p", conn.Id(), conn)
		}),
		ConnClosingHandlerOption(func(ctx context.Context, conn IConnection) {
			log.Info(ctx, "client conn closing: %v, %p", conn.Id(), conn)
		}),
		ConnClosedHandlerOption(func(ctx context.Context, conn IConnection) {
			log.Info(ctx, "client conn closed: %v, %p", conn.Id(), conn)
		}),
		RecvPingHandlerOption(func(ctx context.Context, con IConnection) {
			log.Debug(ctx, "client recv ping")
		}),
		RecvPongHandlerOption(func(ctx context.Context, con IConnection) {
			log.Debug(ctx, "client recv pong")
		}),
	)
	log.Info(ctx, "%v", conn)
	time.Sleep(time.Millisecond * 200)

	for i := 0; i < 10; i++ {
		packet := GetPoolMessage(C2S_REQ)
		packet.SetData([]byte("client rpc req info"))
		resp, err := conn.SendRequestMsg(context.Background(), packet, nil)
		if err != nil {
			t.Error(err)
		}

		if serverResp != string(resp.GetData()) {
			t.Error(resp)
		}

		log.Info(ctx, "client recv: sn: %v, data: %v", resp.GetSn(), string(resp.GetData()))
	}
	time.Sleep(time.Second * 1)
}

func TestWssRequestResponseWithTimeout(t *testing.T) {
	log.SetLogLevel(zapcore.WarnLevel) //if you need info log, commmet this line

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

		time.Sleep(time.Second * 2)
		packet := GetPoolMessage(S2C_RESP_TIMEOUT)
		packet.SetData([]byte("server rpc resp info timeout"))
		return connection.SendResponseMsg(ctx, packet, message.GetSn(), nil)
	})

	//server start
	handler := http.NewServeMux()
	handler.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
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
	go http.ListenAndServe(":8003", handler)

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
	time.Sleep(time.Millisecond * 200)

	toCtx, _ := context.WithTimeout(ctx, time.Second*1)
	packet := GetPoolMessage(C2S_REQ_TIMEOUT)
	packet.SetData([]byte("client rpc req info timeout"))
	resp, err := conn.SendRequestMsg(toCtx, packet, nil)

	if err != ErrWsRpcResponseTimeout {
		t.Error(err)
	}

	if err == nil {
		log.Info(ctx, "client recv: sn: %v, data: %v", resp.GetSn(), string(resp.GetData()))
	} else {
		log.Error(ctx, "client recv err: %v", err)
	}
}

func TestWssSendMessage(t *testing.T) {
	log.SetLogLevel(zapcore.WarnLevel) //if you need info log, commmet this line

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

			log.Debug(ctx, "first pull: %v, %p", conn.Id(), conn)

		}, func(ctx context.Context, pullConn IConnection) {
			log.Debug(ctx, "pull send: %v, %p", conn.Id(), conn)
			//msg from db...
			time.Sleep(time.Millisecond * 100)

			packet := GetPoolMessage(S2C_RESP)
			packet.SetData([]byte("pull msg from db"))
			pullConn.SendMsg(ctx, packet, nil)

			pullConn.SendMsg(ctx, commonMsg, nil)
		})
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
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
				log.Debug(ctx, "server conn establish: %v, %p", conn.Id(), conn)
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
	go http.ListenAndServe(":8003", handler)

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
	time.Sleep(time.Millisecond * 200)

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
	packet := GetPoolMessage(C2S_REQ)
	packet.SetData([]byte("client request2"))
	conn2.SendMsg(context.Background(), packet, nil)

	packet = GetPoolMessage(C2S_REQ)
	packet.SetData([]byte("client request3"))
	conn2.SendMsg(context.Background(), packet, nil)

	conn2.SetCommDataValue("gotuilskey", 100)
	if num, _ := conn2.GetCommDataValue("gotuilskey"); num != 100 {
		t.Error(num)
	}

	conn2.IncrCommDataValueBy("gotuilskey", 100)
	if num, _ := conn2.GetCommDataValue("gotuilskey"); num != 200 {
		t.Error(num)
	}

	conn2.RemoveCommDataValue("gotuilskey")
	if _, ok := conn2.GetCommDataValue("gotuilskey"); ok {
		t.FailNow()
	}
}

func TestWssDialConnect(t *testing.T) {
	log.SetLogLevel(zapcore.WarnLevel) //if you need info log, commmet this line

	InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke
	InitClient()                                                    //client invoke
	ctx := context.Background()

	//server start
	handler := http.NewServeMux()
	handler.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
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
	go http.ListenAndServe(":8003", handler)

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

	time.Sleep(time.Second * 1)

	if len(ClientConnHub.ConnectionIds()) != 100 {
		t.Error(len(ClientConnHub.ConnectionIds()))
	}

	//kick client connect again
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

	time.Sleep(time.Second * 1)

	if len(ClientConnHub.ConnectionIds()) != 100 {
		t.Error(len(ClientConnHub.ConnectionIds()))
	}
}
