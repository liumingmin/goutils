package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/liumingmin/goutils/utils"

	"github.com/gin-gonic/gin"
)

type LimitKeyFunc func(*gin.Context) (string, error)

type LimitReq struct {
	keyFunc LimitKeyFunc

	store map[string]*limitReqRec
	lock  sync.Mutex
}

type limitReqRec struct {
	excess int64
	last   int64
}

func NewLimitReq(keyFunc LimitKeyFunc) *LimitReq {
	lr := &LimitReq{
		keyFunc: keyFunc,
		store:   make(map[string]*limitReqRec),
	}
	return lr
}

func (l *LimitReq) Incoming(keyFunc LimitKeyFunc, rate float64, burst int) gin.HandlerFunc {
	rate = rate * 1000
	burst = burst * 1000

	if keyFunc == nil {
		keyFunc = l.keyFunc
	}

	return func(c *gin.Context) {
		if keyFunc == nil {
			c.Next()
			return
		}

		key, err := keyFunc(c)
		if err != nil {
			c.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}

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
