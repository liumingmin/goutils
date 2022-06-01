package middleware

import (
	"context"
	"net/http"
	"reflect"

	"github.com/liumingmin/goutils/log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap/zapcore"
)

type ServiceFunc func(context.Context, interface{}) (interface{}, error)

type ServiceRespFactory interface {
	IsJsonResponse(data interface{}) bool
	NewErrRespWithCode(code int, err error, data interface{}, tag string) ServiceResponse
	NewDataResponse(data interface{}, tag string) ServiceResponse
}

type ServiceResponse interface {
	GetCode() int
	GetMsg() string
}

const (
	Success   = 0  // 成功
	Unknown   = -1 // 未知错误
	WrongArgs = -2 // 参数错误
)

type Errorx3 interface {
	Error() string
	Code() int
	LogLevel() zapcore.Level
}

func ServiceHandler(serviceFunc ServiceFunc, reqVal interface{}, respFactory ServiceRespFactory) func(*gin.Context) {
	var reqType reflect.Type = nil
	if reqVal != nil {
		value := reflect.Indirect(reflect.ValueOf(reqVal))
		reqType = value.Type()
	}

	if respFactory == nil {
		respFactory = defaultServiceRespFactory
	}

	return func(c *gin.Context) {
		tag := c.Request.RequestURI
		log.Debug(c, tag+" enter")

		var req interface{} = nil
		if reqType != nil {
			req = reflect.New(reqType).Interface()

			if err := c.ShouldBindBodyWith(req, binding.JSON); err != nil {
				log.Error(c, "Bind json failed. error: %v", err)
				c.JSON(http.StatusOK, respFactory.NewErrRespWithCode(WrongArgs, err, nil, tag))
				return
			}
		}

		data, err := serviceFunc(c, req)
		if err != nil {
			lvl := zapcore.ErrorLevel
			code := Unknown
			if errX, ok := err.(Errorx3); ok {
				lvl = errX.LogLevel()
				code = errX.Code()
			}

			if code == Success {
				code = Unknown
			}

			log.Log(c, lvl, tag+" failed. error: %v", err)
			c.JSON(http.StatusOK, respFactory.NewErrRespWithCode(code, err, data, tag))
			return
		}

		if data != nil && !respFactory.IsJsonResponse(data) {
			return
		}

		c.JSON(http.StatusOK, respFactory.NewDataResponse(data, tag))
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var defaultServiceRespFactory = &DefaultServiceRespFactory{}

type DefaultServiceRespFactory struct {
}

func (t *DefaultServiceRespFactory) IsJsonResponse(data interface{}) bool {
	return true
}

func (t *DefaultServiceRespFactory) NewErrRespWithCode(code int, err error, data interface{}, tag string) ServiceResponse {
	var r DefaultServiceResponse
	r.Code = code
	r.Data = data
	r.Tag = tag

	r.Msg = "unknown error"
	if err != nil {
		r.Msg = err.Error()
	}

	return &r
}
func (t *DefaultServiceRespFactory) NewDataResponse(data interface{}, tag string) ServiceResponse {
	var r DefaultServiceResponse
	r.Msg = "success"
	r.Code = Success
	r.Data = data
	r.Tag = tag
	return &r
}

type DefaultServiceResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Tag  string      `json:"tag,omitempty"`
	Data interface{} `json:"data"`
}

func (t *DefaultServiceResponse) GetCode() int {
	return t.Code
}

func (t *DefaultServiceResponse) GetMsg() string {
	return t.Msg
}
