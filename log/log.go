package log

import (
	"context"
	"encoding/base32"

	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

type ILog interface {
	Log(context.Context, zapcore.Level, string, ...interface{})
	Debug(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Recover(context.Context, func(interface{}) string)
	ErrorStack(context.Context, string, ...interface{})

	GetLogLevel() zapcore.Level
	SetLogLevel(zapcore.Level) zapcore.Level
	LogMore() zapcore.Level
	LogLess() zapcore.Level
}

// IFieldsGenerator default json field
type IFieldsGenerator interface {
	Generate(context.Context) []zapcore.Field
	GenerateStr(context.Context) string
}

var LOG_TRACE_CTX_KEY = "__GTraceId__"
var LogImpl = NewZapLogImpl()
var BaseFieldsGenerator IFieldsGenerator = &DefaultFieldsGenerator{Nop: make([]zapcore.Field, 0)}

func Log(ctx context.Context, level zapcore.Level, msg string, args ...interface{}) {
	LogImpl.Log(ctx, level, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...interface{}) {
	LogImpl.Debug(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...interface{}) {
	LogImpl.Info(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...interface{}) {
	LogImpl.Warn(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...interface{}) {
	LogImpl.Error(ctx, msg, args...)
}

func Recover(ctx context.Context, errHandler func(interface{}) string) {
	LogImpl.Recover(ctx, errHandler)
}
func ErrorStack(ctx context.Context, msg string, args ...interface{}) {
	LogImpl.ErrorStack(ctx, msg, args...)
}

func GetLogLevel() zapcore.Level {
	return LogImpl.GetLogLevel()
}

func SetLogLevel(lvl zapcore.Level) zapcore.Level {
	return LogImpl.SetLogLevel(lvl)
}

func LogMore() zapcore.Level {
	return LogImpl.LogMore()
}

func LogLess() zapcore.Level {
	return LogImpl.LogLess()
}

func ContextWithTraceId() context.Context {
	return ContextWithTraceIdFromParent(context.Background())
}

func ContextWithTraceIdFromParent(parent context.Context) context.Context {
	return context.WithValue(parent, LOG_TRACE_CTX_KEY, NewTraceId())
}

func NewTraceId() string {
	uuidBytes := [16]byte(uuid.New())
	b32Uuid := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(uuidBytes[:])
	return b32Uuid
}
