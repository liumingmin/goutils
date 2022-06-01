package middleware

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestServiceHandler(t *testing.T) {
	router := gin.New()
	router.POST("/foo", ServiceHandler(serviceFoo, fooReq{}, nil))

	router.Run(":8080")
}

type fooReq struct {
	FooStr string `json:"fooStr"`
}

func serviceFoo(ctx context.Context, reqVal interface{}) (interface{}, error) {
	req := reqVal.(*fooReq)
	req.FooStr = "respone:" + req.FooStr
	return req, nil
}
