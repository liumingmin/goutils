package errcode

import (
	"go.uber.org/zap/zapcore"
)

const (
	Success   = 0  // 成功
	Unknown   = -1 // 未知错误
	WrongArgs = -2 // 参数错误
)

type Errorx interface {
	Error() string
	Code() int
	LogLevel() zapcore.Level
}

type errorx struct {
	msg      string
	code     int
	logLevel zapcore.Level
}

func (e *errorx) Error() string {
	return e.msg
}

func (e *errorx) Code() int {
	return e.code
}

func (e *errorx) LogLevel() zapcore.Level {
	return e.logLevel
}

func NewErrx(code int, msg string) Errorx {
	return &errorx{
		msg:      msg,
		code:     code,
		logLevel: zapcore.ErrorLevel,
	}
}

func NewErrx2(code int, msg string, logLevel zapcore.Level) Errorx {
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
