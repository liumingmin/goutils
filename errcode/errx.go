package errcode

import (
	"github.com/liumingmin/goutils/log4go"
)

const (
	Success   = 0  // 成功
	Unknown   = -1 // 未知错误
	WrongArgs = -2 // 参数错误
)

type Errorx interface {
	Error() string
	Code() int
	LogLevel() log4go.Level
}

type errorx struct {
	msg      string
	code     int
	logLevel log4go.Level
}

func (e *errorx) Error() string {
	return e.msg
}

func (e *errorx) Code() int {
	return e.code
}

func (e *errorx) LogLevel() log4go.Level {
	return e.logLevel
}

func NewErrx(code int, msg string) Errorx {
	return &errorx{
		msg:      msg,
		code:     code,
		logLevel: log4go.ERROR,
	}
}

func NewErrx2(code int, msg string, logLevel log4go.Level) Errorx {
	return &errorx{
		msg:      msg,
		code:     code,
		logLevel: logLevel,
	}
}

func IsErrorx(e interface{}) bool {
	_, ok := e.(Errorx)
	return ok
}
