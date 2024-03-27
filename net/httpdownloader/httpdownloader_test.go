package httpdownloader

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/liumingmin/goutils/net/bwlimit"
)

var testTempDirPath = filepath.Join(os.TempDir(), "goutils_httpdl")

func TestHttpDownloaderDownload(t *testing.T) {
	os.MkdirAll(testTempDirPath, 0666)

	dialer := bwlimit.NewDialer()
	dialer.RxBwLimit().SetBwLimit(20 * 1024 * 1024)
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
		HttpClient:   hc,
		ConBlockChan: make(chan struct{}, 10),
		BlockSize:    1024 * 1024,
		RetryCnt:     1,
	}

	savePath := filepath.Join(testTempDirPath, "vc_redist")

	url := "https://aka.ms/vs/17/release/vc_redist.arm64.exe"
	resp, err := http.Head(url)
	if err != nil {
		t.Error(err)
	}
	if resp != nil {
		resp.Body.Close()
	}

	contentLen := resp.Header.Get("Content-Length")

	fileSize, _ := strconv.ParseInt(contentLen, 10, 64)

	err = downloader.Download(context.Background(), url, http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Linux; Android 10; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.183 Mobile Safari/537.36"},
	}, savePath)

	if err != nil {
		t.Error(err)
	}

	fileInfo, err := os.Stat(savePath)
	if err != nil {
		t.Error(err)
	}

	if fileInfo.Size() != fileSize {
		t.Error(fileSize)
	}

}

func TestMain(m *testing.M) {

	os.MkdirAll(testTempDirPath, 0666)

	m.Run()

	os.RemoveAll(testTempDirPath)
}
