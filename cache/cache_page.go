package cache

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"goutils/log"
	"goutils/utils"

	"github.com/gin-gonic/gin"
)

func CachePage(store CacheStore, expire time.Duration, reqProc func(r *http.Request) string,
	handle gin.HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		newkey := reqProc(c.Request)
		if bsValue, err := store.Get(newkey); err != nil {
			// replace writer
			log.Debug(context.Background(), "Cache not hit...")
			writer := NewCachedWriterGin(store, expire, c.Writer, newkey)
			c.Writer = writer
			handle(c)
		} else {
			log.Debug(context.Background(), "Cache hit...")
			jsonStr, _ := bsValue.(string)
			var respCache responseCache
			if err := json.Unmarshal([]byte(jsonStr), &respCache); err == nil {
				c.Data(http.StatusOK, "application/json", respCache.Data)
			} else {
				c.String(http.StatusInternalServerError, "cache error")
			}
		}
	}
}

func SimpleRequestFunc(r *http.Request) string {
	url := r.URL
	token := r.Header.Get("token")
	dumpBody, err := utils.DumpBodyAsBytes(r)
	if err != nil {
		log.Error(context.Background(), "Dump request body failed. error: %v", err)
		return ""
	}
	bodyMd5 := utils.MD5(string(dumpBody))

	return utils.UrlEscape("prefix", url.RequestURI()) + token + bodyMd5
}
