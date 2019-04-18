package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	math2 "github.com/liumingmin/goutils/math"
)

type LimitReq struct {
	KeyFunc func(*gin.Context) string
	store   sync.Map
}

type limitReqRec struct {
	excess int64
	last   int64
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

		var excess int64
		nowmill := time.Now().UnixNano() / int64(time.Millisecond)
		if v, ok := l.store.Load(key); ok {
			rec := v.(*limitReqRec)
			elapsed := nowmill - rec.last
			excess = math2.Max64(rec.excess-int64(rate)*math2.Abs64(elapsed)/1000+1000, 0)

			if excess > int64(burst) {
				c.AbortWithStatus(http.StatusServiceUnavailable)
				return
			}
			rec.excess, rec.last = excess, nowmill
		} else {
			rec := &limitReqRec{excess: 0, last: nowmill}
			l.store.Store(key, rec)
		}

		c.Next()
		return
	}
}
