package cache

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liumingmin/goutils/log4go"
)

type responseCache struct {
	Status int
	Header http.Header
	Data   []byte
}

type CachedWriterGin struct {
	gin.ResponseWriter
	status  int
	written bool
	store   CacheStore
	expire  time.Duration
	key     string
}

func NewCachedWriterGin(store CacheStore, expire time.Duration, writer gin.ResponseWriter,
	key string) *CachedWriterGin {
	return &CachedWriterGin{writer, 0, false, store, expire, key}
}

func (w *CachedWriterGin) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *CachedWriterGin) Status() int {
	return w.ResponseWriter.Status()
}

func (w *CachedWriterGin) Written() bool {
	return w.ResponseWriter.Written()
}

func (w *CachedWriterGin) Write(data []byte) (int, error) {
	ret, err := w.ResponseWriter.Write(data)
	if err == nil {
		store := w.store
		if store.Exists(w.key) {
			log4go.Debug("No need cache...")
			return ret, err
		}

		val := responseCache{
			w.status,
			w.Header(),
			data,
		}
		err = store.Set(w.key, val, int64(w.expire/time.Second))
		if err != nil {
			log4go.Error("Cache page store set failed. error: %v", err)
		}
	}
	return ret, err
}

func (w *CachedWriterGin) WriteString(data string) (n int, err error) {
	ret, err := w.ResponseWriter.WriteString(data)
	if err == nil {
		store := w.store
		if store.Exists(w.key) {
			log4go.Debug("No need cache...")
			return ret, err
		}

		//cache response
		val := responseCache{
			w.status,
			w.Header(),
			[]byte(data),
		}

		bs, _ := json.Marshal(val)
		store.Set(w.key, string(bs), int64(w.expire/time.Second))
	}
	return ret, err
}
