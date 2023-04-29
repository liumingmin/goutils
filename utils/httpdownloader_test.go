package utils

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/liumingmin/goutils/net/bwlimit"
)

func TestHttpDownloaderDownload(t *testing.T) {
	dialer := bwlimit.NewDialer()
	dialer.RxBwLimit().SetBwLimit(1024 * 1024)
	hc := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	downloader := &HttpDownloader{
		HttpClient:    hc,
		GoroutinesCnt: 1,
		RetryCnt:      1,
	}

	err := downloader.Download(context.Background(), "xxxxx", http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Linux; Android 10; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Mobile Safari/537.36"},
	}, "xxxx")

	t.Log(err)
}
