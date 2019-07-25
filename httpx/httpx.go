package httpx

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const (
	HTTP_CLIENT_H11 = "HTTP/1.1"
	HTTP_CLIENT_H20 = "HTTP/2.0"
)

type HttpClientX struct {
	Hc11       *http.Client
	Hc20       *http.Client
	protoCache sync.Map
}

func (c *HttpClientX) checkProto(u *url.URL) string {
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

func (c *HttpClientX) Do(req *http.Request) (*http.Response, error) {
	spec := c.checkProto(req.URL)
	if spec == HTTP_CLIENT_H11 {
		return c.Hc11.Do(req)
	} else if spec == HTTP_CLIENT_H20 {
		return c.Hc20.Do(req)
	}
	return nil, errors.New("unknown protocol")
}

func (c *HttpClientX) Get(urlsr string) (resp *http.Response, err error) {
	urlp, err := url.Parse(urlsr)
	if err != nil {
		return nil, err
	}
	spec := c.checkProto(urlp)
	if spec == HTTP_CLIENT_H11 {
		return c.Hc11.Get(urlsr)
	} else if spec == HTTP_CLIENT_H20 {
		return c.Hc20.Get(urlsr)
	}
	return nil, errors.New("unknown protocol")
}

func (c *HttpClientX) Head(urlsr string) (resp *http.Response, err error) {
	urlp, err := url.Parse(urlsr)
	if err != nil {
		return nil, err
	}
	spec := c.checkProto(urlp)
	if spec == HTTP_CLIENT_H11 {
		return c.Hc11.Head(urlsr)
	} else if spec == HTTP_CLIENT_H20 {
		return c.Hc20.Head(urlsr)
	}
	return nil, errors.New("unknown protocol")
}

func (c *HttpClientX) Post(urlsr, contentType string, body io.Reader) (resp *http.Response, err error) {
	urlp, err := url.Parse(urlsr)
	if err != nil {
		return nil, err
	}
	spec := c.checkProto(urlp)
	if spec == HTTP_CLIENT_H11 {
		return c.Hc11.Post(urlsr, contentType, body)
	} else if spec == HTTP_CLIENT_H20 {
		return c.Hc20.Post(urlsr, contentType, body)
	}
	return nil, errors.New("unknown protocol")
}

func (c *HttpClientX) PostForm(urlsr string, data url.Values) (resp *http.Response, err error) {
	urlp, err := url.Parse(urlsr)
	if err != nil {
		return nil, err
	}
	spec := c.checkProto(urlp)
	if spec == HTTP_CLIENT_H11 {
		return c.Hc11.PostForm(urlsr, data)
	} else if spec == HTTP_CLIENT_H20 {
		return c.Hc20.PostForm(urlsr, data)
	}
	return nil, errors.New("unknown protocol")
}
