package middleware

import (
	"net/http"
	"time"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/liumingmin/goutils/utils"
)

type LimitReq struct {
	KeyFunc func(*gin.Context) string
	store   map[string]*limitReqRec
	lock    sync.Mutex
}

type limitReqRec struct {
	excess int64
	last   int64
}

func (l *LimitReq) Init() {
	l.store = make(map[string]*limitReqRec)
}

func (l *LimitReq) Incoming(rate, burst int) gin.HandlerFunc {
	rate = rate * 1000
	burst = burst * 1000

	return func(c *gin.Context) {
		if l.KeyFunc == nil {
			c.Next()
			return
		}

		key := l.KeyFunc(c)
		if key == "" {
			c.Next()
			return
		}

		result := func() bool {
			l.lock.Lock()
			defer l.lock.Unlock()

			var excess int64
			nowmill := time.Now().UnixNano() / int64(time.Millisecond)
			if rec, ok := l.store[key]; ok {
				elapsed := nowmill - rec.last
				excess = utils.Max64(rec.excess-int64(rate)*utils.Abs64(elapsed)/1000+1000, 0)

				if excess > int64(burst) {
					return false
				}
				rec.excess, rec.last = excess, nowmill
			} else {
				rec := &limitReqRec{excess: 0, last: nowmill}
				l.store[key] = rec
			}

			return true
		}()

		if !result {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}

		c.Next()
		return
	}
}
