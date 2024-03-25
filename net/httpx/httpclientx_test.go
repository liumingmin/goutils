package httpx

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func TestHttpXGet(t *testing.T) {
	testRunHttpxServer()

	clientX := getHcx()

	for i := 0; i < 3; i++ {
		resp, err := clientX.Get("http://127.0.0.1:17771/goutils/httpx")
		if err != nil {
			t.Error(fmt.Errorf("error making request: %v", err))
		}

		if resp.StatusCode >= 400 {
			t.Error(resp.StatusCode)
		}

		if resp.Proto != "HTTP/1.1" {
			t.Error(resp.Proto)
		}
	}

	for i := 0; i < 3; i++ {
		resp, err := clientX.Get("http://127.0.0.1:17772/goutils/httpx")
		if err != nil {
			t.Error(fmt.Errorf("error making request: %v", err))
		}

		if resp.StatusCode >= 400 {
			t.Error(resp.StatusCode)
		}

		if resp.Proto != "HTTP/2.0" {
			t.Error(resp.Proto)
		}
	}
}

func TestHttpXPost(t *testing.T) {
	testRunHttpxServer()

	clientX := getHcx()

	for i := 0; i < 3; i++ {
		resp, err := clientX.Post("http://127.0.0.1:17771/goutils/httpx", "application/json", strings.NewReader(""))
		if err != nil {
			t.Error(fmt.Errorf("error making request: %v", err))
		}

		if resp.StatusCode >= 400 {
			t.Error(resp.StatusCode)
		}

		if resp.Proto != "HTTP/1.1" {
			t.Error(resp.Proto)
		}
	}

	for i := 0; i < 3; i++ {
		resp, err := clientX.Post("http://127.0.0.1:17772/goutils/httpx", "application/json", strings.NewReader(""))
		if err != nil {
			t.Error(fmt.Errorf("error making request: %v", err))
		}

		if resp.StatusCode >= 400 {
			t.Error(resp.StatusCode)
		}

		if resp.Proto != "HTTP/2.0" {
			t.Error(resp.Proto)
		}
	}
}

func TestHttpXHead(t *testing.T) {
	testRunHttpxServer()

	clientX := getHcx()

	for i := 0; i < 3; i++ {
		resp, err := clientX.Head("http://127.0.0.1:17771/goutils/httpx")
		if err != nil {
			t.Error(fmt.Errorf("error making request: %v", err))
		}

		if resp.StatusCode >= 400 {
			t.Error(resp.StatusCode)
		}

		if resp.Proto != "HTTP/1.1" {
			t.Error(resp.Proto)
		}
	}

	for i := 0; i < 3; i++ {
		resp, err := clientX.Head("http://127.0.0.1:17772/goutils/httpx")
		if err != nil {
			t.Error(fmt.Errorf("error making request: %v", err))
		}

		if resp.StatusCode >= 400 {
			t.Error(resp.StatusCode)
		}

		if resp.Proto != "HTTP/2.0" {
			t.Error(resp.Proto)
		}
	}
}

func TestHttpXPostForm(t *testing.T) {
	testRunHttpxServer()

	clientX := getHcx()

	for i := 0; i < 3; i++ {
		resp, err := clientX.PostForm("http://127.0.0.1:17771/goutils/httpx", url.Values{"key": []string{"value"}})
		if err != nil {
			t.Error(fmt.Errorf("error making request: %v", err))
		}

		if resp.StatusCode >= 400 {
			t.Error(resp.StatusCode)
		}

		if resp.Proto != "HTTP/1.1" {
			t.Error(resp.Proto)
		}
	}

	for i := 0; i < 3; i++ {
		resp, err := clientX.PostForm("http://127.0.0.1:17772/goutils/httpx", url.Values{"key": []string{"value"}})
		if err != nil {
			t.Error(fmt.Errorf("error making request: %v", err))
		}

		if resp.StatusCode >= 400 {
			t.Error(resp.StatusCode)
		}

		if resp.Proto != "HTTP/2.0" {
			t.Error(resp.Proto)
		}
	}
}

func getHcx() *HttpClientX {
	return &HttpClientX{
		Hc11: &http.Client{},
		Hc20: &http.Client{
			// Skip TLS dial
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
					return net.Dial(network, addr)
				},
			},
		},
	}
}

func BenchmarkHttpx(b *testing.B) {
	testRunHttpxServer()
	clientX := getHcx()
	//b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := clientX.Get("http://127.0.0.1:17772/goutils/httpx")
		if err != nil {
			b.Fatal(fmt.Errorf("error making request: %v", err))
		}
		//b.Log(resp.StatusCode)
		//b.Log(resp.Proto)
	}
}

func testRunHttpxServer() {
	handler := http.NewServeMux()
	handler.HandleFunc("/goutils/httpx", func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		w.Write(data)
	})

	go http.ListenAndServe(":17771", handler)

	h2s := &http2.Server{}
	h1s := &http.Server{
		Addr:    ":17772",
		Handler: h2c.NewHandler(handler, h2s),
	}
	go h1s.ListenAndServe()
}
