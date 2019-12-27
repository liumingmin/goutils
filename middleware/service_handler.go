package middleware

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/liumingmin/goutils/errcode"
	"github.com/liumingmin/goutils/log4go"
	"github.com/liumingmin/goutils/model"
)

type ServiceFunc func(context.Context, interface{}) (interface{}, int, error)

func ServiceHandler(serviceFunc ServiceFunc, reqVal interface{}) func(*gin.Context) {
	var reqType reflect.Type = nil
	if reqVal != nil {
		value := reflect.Indirect(reflect.ValueOf(reqVal))
		reqType = value.Type()
	}

	return func(c *gin.Context) {
		tag := c.Request.RequestURI
		log4go.Info(tag + " enter")

		var req interface{} = nil
		if reqType != nil {
			req = reflect.New(reqType).Interface()

			//使用http 200 ok 响应code
			if err := c.ShouldBindWith(req, binding.JSON); err != nil {
				log4go.Error("Bind json failed. error: %v", err)
				c.JSON(http.StatusOK, model.NewBindFailedResponse(tag))
				return
			}
		}

		resp, code, err := serviceFunc(c, req)
		if err != nil {
			lvl := log4go.ERROR
			if code == 0 {
				if err828x, ok := err.(errcode.Errorx); ok {
					code = err828x.Code()
					lvl = err828x.LogLevel()
				}

				if code == 0 {
					code = errcode.Unknown
				}
			}

			if lvl == log4go.ERROR {
				log4go.Error(tag+" failed. error: %v", err)
			} else {
				log4go.Info(tag+" failed. error: %v", err)
			}

			c.JSON(http.StatusOK, model.NewErrRespWithCode(code, err, tag))
			return
		}
		c.JSON(http.StatusOK, model.NewDataResponse(resp, tag))
	}
}
