package ws

import (
	"time"

	"github.com/liumingmin/goutils/conf"
)

//配置项
var (
	dispatcherNum   = conf.ExtInt("ws.dispatcherNum", 16)              //并发处理消息数量
	maxFailureRetry = conf.ExtInt("ws.maxFailureRetry", 10)            //重试次数
	ReadWait        = conf.ExtDuration("ws.readWait", 60*time.Second)  //读等待
	WriteWait       = conf.ExtDuration("ws.writeWait", 60*time.Second) //写等待
	PingPeriod      = WriteWait * 4 / 10                               //ping间隔应该小于写等待时间

	NetTemporaryWait = 500 * time.Millisecond //网络抖动重试等待
)

//客户端独有配置项
var (
	handshakeTimeout = conf.ExtDuration("ws.dialTimeout", "10s")
	connMaxRetry     = conf.ExtInt("ws.connMaxRetry", 10)
)
