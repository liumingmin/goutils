package ws

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/demdxx/gocast"
	"github.com/gorilla/websocket"
	"github.com/liumingmin/goutils/log"
	"github.com/liumingmin/goutils/utils/safego"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
)

func TestWsPb(t *testing.T) {
	if *P_BASE_s2c_err_displace.Enum() != P_BASE_s2c_err_displace {
		t.FailNow()
	}

	if int32(P_BASE_s2c_err_displace.Number()) != int32(P_BASE_s2c_err_displace) {
		t.FailNow()
	}
	if P_BASE_s2c_err_displace.String() != "s2c_err_displace" {
		t.Error(P_BASE_s2c_err_displace.String())
	}

	displace := &P_DISPLACE{}
	displace.OldIp = []byte("10.1.1.1")
	displace.NewIp = []byte("10.1.1.2")
	displace.Ts = time.Now().UnixNano()

	bs, err := proto.Marshal(displace)
	if err != nil {
		t.Error(err)
	}

	displace2 := &P_DISPLACE{}
	err = proto.Unmarshal(bs, displace2)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(displace.GetNewIp(), displace2.NewIp) {
		t.FailNow()
	}

	if !bytes.Equal(displace.GetOldIp(), displace2.OldIp) {
		t.FailNow()
	}

	if displace.GetTs() != displace2.Ts {
		t.FailNow()
	}
}

func TestWsHub(t *testing.T) {
	InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke

	conn, err := ClientConnHub.Find("dummy")
	if err == nil {
		t.Error(conn)
	}
}

func TestWsOption(t *testing.T) {
	upgrader := &websocket.Upgrader{}
	conn := &Connection{}
	SrvUpgraderOption(upgrader)(conn)

	if conn.upgrader != upgrader {
		t.Error(conn.upgrader)
	}

	pullChannelIds := []int{}
	SrvPullChannelsOption(pullChannelIds)(conn)
	if conn.pullChannelMap != nil {
		t.Error(conn.pullChannelMap)
	}

	pullChannelIds = []int{0, 1, 2}
	SrvPullChannelsOption(pullChannelIds)(conn)
	if len(conn.pullChannelMap) != 3 {
		t.Error(conn.pullChannelMap)
	}

	ClientDialOption(&websocket.Dialer{})(conn)

	ClientDialWssOption("gou://127.0.0.1:8080", false)(conn)
	if conn.dialer.TLSClientConfig != nil {
		t.Error(conn.dialer.TLSClientConfig)
	}

	ClientDialWssOption("wss://127.0.0.1:8080", false)(conn)

	if !conn.dialer.TLSClientConfig.InsecureSkipVerify {
		t.FailNow()
	}

	NetMaxFailureRetryOption(-1)(conn)
	if conn.maxFailureRetry == -1 {
		t.Error(conn.maxFailureRetry)
	}

	NetReadWaitOption(-1)(conn)
	if conn.readWait == -1 {
		t.Error(conn.readWait)
	}

	NetWriteWaitOption(-1)(conn)
	if conn.writeWait == -1 {
		t.Error(conn.writeWait)
	}

	NetTemporaryWaitOption(-1)(conn)
	if conn.temporaryWait == -1 {
		t.Error(conn.temporaryWait)
	}
}

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
			Source:   "",
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
			Source:   "",
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

	uid := "100"
	typed := 2
	deviceId := "a100"
	version := 2
	charset := 1

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
			Typed:    gocast.ToInt(r.URL.Query().Get("typed")),
			DeviceId: r.URL.Query().Get("deviceId"),
			Source:   r.URL.Query().Get("source"),
			Version:  gocast.ToInt(r.URL.Query().Get("version")),
			Charset:  gocast.ToInt(r.URL.Query().Get("charset")),
		}
		_, err := Accept(ctx, w, r, connMeta, DebugOption(true),
			SrvUpgraderCompressOption(true),
			CompressionLevelOption(2),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				if conn.(*Connection).UserId() != uid {
					t.Error(conn.UserId())
				}

				if conn.(*Connection).Type() != typed {
					t.Error(conn.Type())
				}

				if conn.(*Connection).DeviceId() != deviceId {
					t.Error(conn.DeviceId())
				}

				if conn.(*Connection).Version() != version {
					t.Error(conn.Version())
				}

				if conn.(*Connection).Charset() != charset {
					t.Error(conn.Charset())
				}

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

	url := fmt.Sprintf("ws://127.0.0.1:8003/join?uid=%v&typed=%v&deviceId=%v&version=%v&charset=%v", uid, typed, deviceId, version, charset)
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
			Source:   "",
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

	dummmyMeta := ConnectionMeta{UserId: "a0"}

	conn0, err := ClientConnHub.Find(dummmyMeta.BuildConnId())
	if err != nil {
		t.Error(err)
	}

	if conn0.UserId() != "a0" {
		t.Error(conn0)
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

func TestWssDialRetryConnect(t *testing.T) {
	log.SetLogLevel(zapcore.WarnLevel) //if you need info log, commmet this line

	InitServerWithOpt(ServerOption{[]HubOption{HubShardOption(4)}}) //server invoke
	InitClient()                                                    //client invoke
	ctx := context.Background()

	var dialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: time.Millisecond * 500,
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
	}

	url := "ws://127.0.0.1:8003/join?uid=a1"
	go func() {
		DialConnect(context.Background(), url, http.Header{},
			DebugOption(true),
			ClientIdOption("a1"),
			ClientDialOption(dialer),
			ClientDialWssOption(url, false),
			ClientDialRetryOption(4, time.Millisecond*100),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "client conn establish: %v", conn.Id())
			}),
		)
	}()

	time.Sleep(time.Millisecond * 500)

	//server start
	handler := http.NewServeMux()
	handler.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		connMeta := ConnectionMeta{
			UserId:   r.URL.Query().Get("uid"),
			Typed:    0,
			DeviceId: "",
			Source:   "",
			Version:  0,
			Charset:  0,
		}
		_, err := Accept(ctx, w, r, connMeta, DebugOption(true),
			SrvUpgraderCompressOption(true),
			CompressionLevelOption(2),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn establish: %v, %p", conn.Id(), conn)
				if len(ClientConnHub.ConnectionIds()) != 1 {
					t.Error(len(ClientConnHub.ConnectionIds()))
				}
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
	time.Sleep(time.Second)

}

func TestAutoReDialConnect(t *testing.T) {
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
			Source:   "",
			Version:  0,
			Charset:  0,
		}
		_, err := Accept(ctx, w, r, connMeta, DebugOption(true),
			SrvUpgraderCompressOption(true),
			CompressionLevelOption(2),
			ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
				log.Info(ctx, "server conn establish: %v, %p", conn.Id(), conn)

				safego.Go(func() {
					time.Sleep(time.Millisecond * 100)
					conn.KickClient(false)
				})
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

	time.Sleep(time.Millisecond * 500)

	url := "ws://127.0.0.1:8003/join?uid=a1"

	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	AutoReDialConnect(ctxTimeout, url, http.Header{}, time.Millisecond*500,
		DebugOption(true),
		ClientDialWssOption(url, false),
		//ClientDialRetryOption(4, time.Millisecond*100),
		ConnEstablishHandlerOption(func(ctx context.Context, conn IConnection) {
			log.Info(ctx, "client conn establish: %v, %p", conn.Id(), conn)
		}),
	)
}

func TestDefaultPuller(t *testing.T) {
	firstPull := false
	pullSendCnt := 0

	conn := &Connection{}
	conn.init()

	pullChannelIds := []int{1}
	SrvPullChannelsOption(pullChannelIds)(conn)

	puller := NewDefaultPuller(conn, 1, func(ctx context.Context, i IConnection) {
		firstPull = true
	}, func(ctx context.Context, i IConnection) {
		pullSendCnt++
	})

	go func() {
		puller.PullSend()
	}()

	go func() {
		time.Sleep(time.Millisecond * 100)

		if puller.(*defaultPuller).isRunning != 1 {
			t.Error(puller.(*defaultPuller).isRunning)
		}

		puller.PullSend()
	}()

	time.Sleep(time.Millisecond * 100)

	if !firstPull {
		t.Error(firstPull)
	}

	if pullSendCnt != 1 {
		t.Error(pullSendCnt)
	}

	conn.SignalPullSend(context.Background(), 1)
	time.Sleep(time.Millisecond * 100)
	if pullSendCnt != 2 {
		t.Error(pullSendCnt)
	}

	conn.SignalPullSend(context.Background(), 2)
	time.Sleep(time.Millisecond * 100)
	if pullSendCnt != 2 {
		t.Error(pullSendCnt)
	}

	firstPullNo := false
	pullSendNo := false
	pullerNoChan := NewDefaultPuller(conn, 2, func(ctx context.Context, i IConnection) {
		firstPullNo = true
	}, func(ctx context.Context, i IConnection) {
		pullSendNo = true
	})

	go func() {
		pullerNoChan.PullSend()
	}()
	time.Sleep(time.Millisecond * 100)
	if firstPullNo {
		t.Error(firstPullNo)
	}

	if pullSendNo {
		t.Error(pullSendNo)
	}

	conn.SignalPullSend(context.Background(), 1)
	conn.closePull(context.Background())

	chann, _ := conn.GetPullChannel(1)
	_, ok := <-chann
	if ok {
		t.Error(chann)
	}
}
