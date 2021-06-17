package middleware

import (
	"context"
	"net/http"
	"reflect"

	"github.com/liumingmin/goutils/errcode"
	"github.com/liumingmin/goutils/log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap/zapcore"
)

type ServiceFunc func(context.Context, interface{}) (interface{}, error)

type ServiceResponse interface {
	IsJsonResponse(data interface{}) bool
	NewBindFailedResponse(err error, tag string) interface{}
	NewErrRespWithCode(code int, err error, data interface{}, tag string) interface{}
	NewDataResponse(data interface{}, tag string) interface{}
}

func ServiceHandler(serviceFunc ServiceFunc, reqVal interface{}, sResp ServiceResponse) func(*gin.Context) {
	var reqType reflect.Type = nil
	if reqVal != nil {
		value := reflect.Indirect(reflect.ValueOf(reqVal))
		reqType = value.Type()
	}

	return func(c *gin.Context) {
		tag := c.Request.RequestURI
		log.Debug(c, tag+" enter")

		var req interface{} = nil
		if reqType != nil {
			req = reflect.New(reqType).Interface()

			if err := c.ShouldBindBodyWith(req, binding.JSON); err != nil {
				log.Error(c, "Bind json failed. error: %v", err)
				c.JSON(http.StatusOK, sResp.NewBindFailedResponse(err, tag))
				return
			}
		}

		data, err := serviceFunc(c, req)
		if err != nil {
			lvl := zapcore.ErrorLevel
			code := errcode.Unknown
			if errX, ok := err.(errcode.Errorx); ok {
				lvl = errX.LogLevel()
				code = errX.Code()
			}

			if code == errcode.Success {
				code = errcode.Unknown
			}

			if lvl == zapcore.ErrorLevel {
				log.Error(c, tag+" failed. error: %v", err)
			} else {
				log.Info(c, tag+" failed. error: %v", err)
			}

			c.JSON(http.StatusOK, sResp.NewErrRespWithCode(code, err, data, tag))
			return
		}

		if data != nil && !sResp.IsJsonResponse(data) {
			return
		}

		c.JSON(http.StatusOK, sResp.NewDataResponse(data, tag))
	}
}
