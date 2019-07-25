package httpx

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"testing"

	"golang.org/x/net/http2"
)

func TestHttpXGet(t *testing.T) {
	clientX := getHcx()

	for i := 0; i < 3; i++ {
		resp, err := clientX.Get("http://120.92.169.81:8049") //http://10.11.253.5:8080
		if err != nil {
			t.Fatal(fmt.Errorf("error making request: %v", err))
		}
		t.Log(resp.StatusCode)
		t.Log(resp.Proto)
	}
}

func TestHttpXPost(t *testing.T) {
	clientX := getHcx()

	for i := 0; i < 3; i++ {
		resp, err := clientX.Get("http://120.92.169.81:8881") //http://10.11.253.5:8080
		if err != nil {
			t.Fatal(fmt.Errorf("error making request: %v", err))
		}
		t.Log(resp.StatusCode)
		t.Log(resp.Proto)
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
	clientX := getHcx()
	//b.ResetTimer()
	b.N = 1000
	for i := 0; i < b.N; i++ {
		_, err := clientX.Get("http://120.92.169.81:8049") //http://10.11.253.5:8080
		if err != nil {
			b.Fatal(fmt.Errorf("error making request: %v", err))
		}
		//b.Log(resp.StatusCode)
		//b.Log(resp.Proto)
	}
}
