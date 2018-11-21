package cbk

import (
	"strings"
	"github.com/gin-gonic/gin"
	"bytes"
	"sync"
	"github.com/liumingmin/goutils/lighttimer"
	"time"
	"net/http"
)

var (
	reqCountMap  = make(map[string]uint32)
	reqLastTimeMap  = make(map[string]time.Time)
	reqCountMapMutex sync.Mutex

	reqBlockMap sync.Map

	lightTimer *lighttimer.LightTimer
)

type Options struct{
	MaxQps   uint32
	CheckSecond uint32
	RecoverSecond uint32
	ReqTagFunc  func(c *gin.Context)(string)
}

func cbUri(c *gin.Context, keyValue string) string {
	method := c.Request.Method
	path := c.Request.URL.Path

	endchar := strings.Index(path,"?")
	if endchar<0 {
		endchar = len(path)
	}

	reqPath :=path[0:endchar]

	sb := bytes.Buffer{}
	sb.WriteString(method)
	sb.WriteString("-")
	sb.WriteString(reqPath)

	if len(keyValue) > 0{
		sb.WriteString("-")
		sb.WriteString(keyValue)
	}

	result := sb.String()

	//fmt.Println(result)
	return result
}

func CircuitBreaker(options Options) gin.HandlerFunc {
	if lightTimer == nil{
		lightTimer =lighttimer.NewLightTimer()
		lightTimer.StartTicks(time.Millisecond)
	}

	if options.CheckSecond == 0{
		options.CheckSecond = 1
	}

	if options.RecoverSecond == 0{
		options.RecoverSecond = 5
	}

	return func(c *gin.Context) {
		cbDone := c.GetHeader("__cb_done__")
		if len(cbDone)>0{
			c.Next()
			return
		}

		defer c.Header("__cb_done__","done")

		tag :=""
		if options.ReqTagFunc != nil{
			tag = options.ReqTagFunc(c)
		}

		cburi := cbUri(c,tag)

		if blocked,isok := reqBlockMap.Load(cburi);isok{
			if blocked.(bool) {
				c.String(http.StatusServiceUnavailable,"To many requests in a second")
				c.Abort()
				return
			}
		}

		//fmt.Println(cburi)

		reqCountMapMutex.Lock()
		if count,isok := reqCountMap[cburi];isok{
			count = count+1
			reqCountMap[cburi]=count

			checkIsBlocked(cburi,count,options)
		}else{
			reqCountMap[cburi]=1
			reqLastTimeMap[cburi]=time.Now()
		}
		reqCountMapMutex.Unlock()

		c.Next()
		return
	}

}

func checkIsBlocked(cburi string, count uint32, options Options)  {
	timeinterval := uint32(time.Since(reqLastTimeMap[cburi])/time.Second)
	if timeinterval > options.CheckSecond{
		if  count/timeinterval >options.MaxQps{
			//fmt.Println(count/timeinterval)
			reqBlockMap.Store(cburi,true)

			lightTimer.AddCallback(time.Second*time.Duration(options.RecoverSecond), func() {
				reqCountMapMutex.Lock()
				delete(reqCountMap, cburi)
				delete(reqLastTimeMap,cburi)
				reqCountMapMutex.Unlock()

				reqBlockMap.Store(cburi,false)
			})
		}else{
			reqCountMap[cburi]=1
			reqLastTimeMap[cburi]=time.Now()
		}
	}
}

