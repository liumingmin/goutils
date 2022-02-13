package middleware

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type LimitConn struct {
	keyFunc LimitKeyFunc
	store   map[string]int
	lock    sync.Mutex
}

func NewLimitConn(keyFunc LimitKeyFunc) *LimitConn {
	lr := &LimitConn{
		keyFunc: keyFunc,
		store:   make(map[string]int),
	}
	return lr
}

func (l *LimitConn) Incoming(keyFunc LimitKeyFunc, max, burst int) gin.HandlerFunc {
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

func (l *LimitConn) Leaving(keyFunc LimitKeyFunc) gin.HandlerFunc {
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

func reqHostIp(c *gin.Context) (string, error) {
	return strings.Split(c.Request.RemoteAddr, ":")[0], nil
}
