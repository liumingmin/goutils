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
	e := gin.Default()
	e.GET("/join", join)
	go e.Run(":8003")

	connectWss("100")
	time.Sleep(time.Minute * 5)
}

func connectWss(uid string) {
	conn, _ := Connect(context.Background(), "server1", "ws://127.0.0.1:8003/join?uid="+uid, false, http.Header{})
	go func() {
		time.Sleep(time.Minute * 2)
		conn.KickServer(false)
	}()
}

func join(ctx *gin.Context) {
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
	//go func() {
	//	time.Sleep(time.Minute * 2)
	//	con.KickClient(false)
	//}()
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
