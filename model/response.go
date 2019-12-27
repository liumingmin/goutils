package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liumingmin/goutils/errcode"
	"github.com/liumingmin/goutils/log4go"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Tag  string      `json:"tag,omitempty"`
	Data interface{} `json:"data"`
}

func NewResponse(err error, args ...interface{}) *Response {
	var data interface{}
	var tag string
	switch len(args) {
	case 1:
		tag = args[0].(string)
	case 2:
		tag = args[0].(string)
		data = args[1]
	}

	if err == nil {
		return NewDataResponse(data, tag)
	} else {
		return NewErrorResponse(err, tag)
	}
}

func NewSuccessResponse(tags ...string) *Response {
	if len(tags) <= 0 {
		return &Response{errcode.Success, "success", "", nil}
	}
	return &Response{errcode.Success, "success", tags[0], nil}
}

func NewErrRespWithCode(code int, err error, args ...interface{}) *Response {
	var r Response
	r.Msg = "unknown error"
	if err != nil && errcode.IsErrorx(err) {
		r.Msg = err.Error()
	}
	r.Code = code
	switch len(args) {
	case 1:
		r.Tag = args[0].(string)
	case 2:
		r.Tag = args[0].(string)
		r.Data = args[1]
	}
	return &r
}

func NewErrorResponse(err error, tag ...interface{}) *Response {
	return NewErrRespWithCode(errcode.Unknown, err, tag...)
}

func NewDataResponse(data interface{}, tag ...string) *Response {
	var r Response
	r.Msg = "success"
	r.Code = errcode.Success
	r.Data = data
	if len(tag) > 0 {
		r.Tag = tag[0]
	}
	return &r
}

func SimpleResponse(code int, msg string) *Response {
	var r Response
	r.Msg = msg
	r.Code = code
	return &r
}

func SimpleSuccessResponse(msg string) *Response {
	var r Response
	r.Msg = msg
	r.Code = errcode.Success
	return &r
}

func NewBindFailedResponse(tag string) *Response {
	return &Response{Code: errcode.WrongArgs, Msg: "wrong argument", Tag: tag}
}

func HandleUnknownSystemError(c *gin.Context) {
	err := recover()
	if err != nil {
		log4go.Critical(err, "unexpected error occurred")
		c.JSON(http.StatusOK, SimpleResponse(errcode.Unknown, "Unknown system error"))
	}
}
