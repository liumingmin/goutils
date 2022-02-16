package log

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestZap(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set("__traceId", "aaabbbbbcccc")
	//Info(ctx, "我是日志", "name", "管理员")  //json

	Info(ctx, "我是日志2")

	//Info(ctx, "我是日志3", "name")  //json

	Error(ctx, "我是日志4: %v,%v", "管理员", "eee")
}

func TestZapJson(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set("__traceId", "aaabbbbbcccc")
	Info(ctx, "我是日志 %v", "name", "管理员") //json

	Info(ctx, "我是日志3", "管理员") //json
	Error(ctx, "我是日志3")       //json
	Log(ctx, zapcore.ErrorLevel, "日志啊")
}

func TestPanicLog(t *testing.T) {
	testPanicLog()

	Info(context.Background(), "catch panic")
}

func testPanicLog() {
	ctx := &gin.Context{}
	ctx.Set("__traceId", "aaabbbbbcccc")

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
