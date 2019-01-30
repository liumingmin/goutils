package cache

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liumingmin/goutils/log4go"
	"github.com/liumingmin/goutils/utils"
)

func CachePage(store CacheStore, expire time.Duration, reqProc func(r *http.Request) string,
	handle gin.HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		newkey := reqProc(c.Request)
		if cache, err := store.Get(newkey); err != nil {
			// replace writer
			log4go.Debug("Cache not hit...")
			writer := NewCachedWriterGin(store, expire, c.Writer, newkey)
			c.Writer = writer
			handle(c)
		} else {
			respCache := cache.(responseCache)
			c.JSON(http.StatusOK, respCache.Data)
		}
	}
}

func SimpleRequestFunc(r *http.Request) string {
	url := r.URL
	token := r.Header.Get("token")
	dumpBody, err := utils.DumpBodyAsBytes(r)
	if err != nil {
		log4go.Error("Dump request body failed. error: %v", err)
		return ""
	}
	bodyMd5 := utils.MD5(string(dumpBody))

	return utils.UrlEscape("prefix", url.RequestURI()) + token + bodyMd5
}
