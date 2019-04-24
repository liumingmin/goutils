package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type LimitConn struct {
	KeyFunc func(*gin.Context) string
	store   map[string]int
	lock    sync.Mutex
}

func (l *LimitConn) Init() {
	l.store = make(map[string]int)
}

func (l *LimitConn) Incoming(max, burst int) gin.HandlerFunc {

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

			if conn, ok := l.store[key]; ok {
				if conn+1 > max+burst {
					return false
				} else {
					l.store[key] = conn + 1
				}
			} else {
				l.store[key] = 1
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

func (l *LimitConn) Leaving() gin.HandlerFunc {
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

		func() {
			l.lock.Lock()
			defer l.lock.Unlock()

			if conn, ok := l.store[key]; ok {
				l.store[key] = conn - 1
			}
		}()

		c.Next()
		return
	}
}
