package log

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type GameDefaultFieldGenerator struct {
}

func (f *GameDefaultFieldGenerator) GetDefaultFields() []zap.Field {
	return []zap.Field{
		zap.String("gameCode", "lol"),
		zap.String("version", "1.0"),
	}
}

func TestZap(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set(LOG_TRADE_ID, "aaabbbbbcccc")

	Info(ctx, "我是日志2")
	SetDefaultGenerator(new(GameDefaultFieldGenerator))
	Error(ctx, "我是日志4: %v,%v", "管理员", "eee")
}

func TestErrorStack(t *testing.T) {
	ErrorStack(context.Background(), "panic error")
}

func TestPanicLog(t *testing.T) {
	testPanicLog()
	Info(context.Background(), "catch panic")
}

func testPanicLog() {
	ctx := &gin.Context{}
	ctx.Set(LOG_TRADE_ID, "aaabbbbbcccc")

	defer Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("recover from err: %v", e)
	})

	panic(errors.New("dddd"))
}

func TestLevelChange(t *testing.T) {
	traceId := time.Now().Unix()
	ctx := context.WithValue(context.Background(), LOG_TRADE_ID, traceId)
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())
	Error(ctx, LogLess())

	fmt.Println(LogLess(), "============")

	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())
	Info(ctx, LogMore())

	fmt.Println(LogMore(), "============")
}
