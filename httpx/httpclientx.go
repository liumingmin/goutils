package httpx

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	HTTP_CLIENT_H10 = "HTTP/1.0"
	HTTP_CLIENT_H11 = "HTTP/1.1"
	HTTP_CLIENT_H20 = "HTTP/2.0"
)

type HttpClientX struct {
	Hc11       *http.Client
	Hc20       *http.Client
	protoCache sync.Map
}

func (c *HttpClientX) checkProto(u *url.URL) string {
	if strings.HasPrefix(strings.ToLower(u.Scheme), "https") {
		return HTTP_CLIENT_H11
	}

	baseUrl := fmt.Sprintf("%s://%s", u.Scheme, u.Host) //[scheme:][//[userinfo@]host]
	if v, ok := c.protoCache.Load(baseUrl); ok {
		//fmt.Println("hit cache:", baseUrl)
		return v.(string)
	} else {
		resp, err := c.Hc20.Head(baseUrl)
		if err == nil && resp != nil {
			c.protoCache.Store(baseUrl, resp.Proto)
			return resp.Proto
		}

		resp, err = c.Hc11.Head(baseUrl)
		if err == nil && resp != nil {
			c.protoCache.Store(baseUrl, resp.Proto)
			return resp.Proto
		}
	}
	return ""
}

func (c *HttpClientX) getClient(proto string) *http.Client {
	if proto == HTTP_CLIENT_H11 || proto == HTTP_CLIENT_H10 {
		return c.Hc11
	} else if proto == HTTP_CLIENT_H20 {
		return c.Hc20
	}
	return nil
}

func (c *HttpClientX) Do(req *http.Request) (*http.Response, error) {
	proto := c.checkProto(req.URL)
	if proto == HTTP_CLIENT_H20 && req != nil {
		req.Header.Del("Connection")
	}

	client := c.getClient(proto)
	if client != nil {
		return client.Do(req)
	}

	return nil, errors.New("unknown protocol")
}

func (c *HttpClientX) Get(urlsr string) (resp *http.Response, err error) {
	urlp, err := url.Parse(urlsr)
	if err != nil {
		return nil, err
	}
	proto := c.checkProto(urlp)
	client := c.getClient(proto)
	if client != nil {
		return client.Get(urlsr)
	}
	return nil, errors.New("unknown protocol")
}

func (c *HttpClientX) Head(urlsr string) (resp *http.Response, err error) {
	urlp, err := url.Parse(urlsr)
	if err != nil {
		return nil, err
	}
	proto := c.checkProto(urlp)
	client := c.getClient(proto)
	if client != nil {
		return client.Head(urlsr)
	}

	return nil, errors.New("unknown protocol")
}

func (c *HttpClientX) Post(urlsr, contentType string, body io.Reader) (resp *http.Response, err error) {
	urlp, err := url.Parse(urlsr)
	if err != nil {
		return nil, err
	}
	proto := c.checkProto(urlp)
	client := c.getClient(proto)
	if client != nil {
		return client.Post(urlsr, contentType, body)
	}
	return nil, errors.New("unknown protocol")
}

func (c *HttpClientX) PostForm(urlsr string, data url.Values) (resp *http.Response, err error) {
	urlp, err := url.Parse(urlsr)
	if err != nil {
		return nil, err
	}
	proto := c.checkProto(urlp)
	client := c.getClient(proto)
	if client != nil {
		return client.PostForm(urlsr, data)
	}

	return nil, errors.New("unknown protocol")
}
