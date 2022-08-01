package log

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestZap(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set(LOG_TRADE_ID, "aaabbbbbcccc")

	Info(ctx, "我是日志2")
	Error(ctx, "我是日志4: %v,%v", "管理员", "eee")
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
	traceId := strings.Replace(uuid.New().String(), "-", "", -1)
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
