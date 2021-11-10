package utils

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var gTransport *http.Client

func UrlEscape(prefix string, u string) string {
	key := url.QueryEscape(u)
	if len(key) > 200 {
		h := sha1.New()
		io.WriteString(h, u)
		key = string(h.Sum(nil))
	}
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(":")
	buffer.WriteString(key)
	return buffer.String()
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func DumpBodyAsReader(req *http.Request) (reader io.ReadCloser, err error) {
	if req == nil || req.Body == nil {
		return nil, errors.New("request or body is nil")
	} else {
		reader, req.Body, err = drainBody(req.Body)
	}
	return
}

func DumpBodyAsBytes(req *http.Request) (copy []byte, err error) {
	var reader io.ReadCloser
	reader, err = DumpBodyAsReader(req)
	copy, err = ioutil.ReadAll(reader)
	return
}

func defaultPooledTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          10240,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		MaxIdleConnsPerHost:   2048,
		//MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
	return transport
}

func DefaultPooledClient() *http.Client {
	return gTransport
}

func PerformTestRequest(method, target string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}

func init() {
	gTransport = &http.Client{
		Transport: defaultPooledTransport(),
	}
}

func ReqHostIp(c *gin.Context) (string, error) {
	return strings.Split(c.Request.RemoteAddr, ":")[0], nil
}

func CopyHttpHeader(from http.Header) http.Header {
	r := http.Header{}
	if from == nil {
		return r
	}

	for k := range from {
		r.Set(k, from.Get(k))
	}
	return r
}
