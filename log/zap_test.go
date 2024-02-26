package log

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go.uber.org/zap"
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
	ctx := context.WithValue(context.Background(), LOG_TRACE_ID, "zap_trace_id")
	Info(ctx, "我是日志2")
	SetDefaultGenerator(new(GameDefaultFieldGenerator))
	Error(ctx, "我是日志4: %v,%v", "管理员", "eee")

	Info(ctx, "我是日志5: %v", "hello")
	Warn(ctx, "我是日志6: %v", "hello sam")
}

func TestErrorStack(t *testing.T) {
	ErrorStack(context.Background(), "panic error")
}

func TestPanicLog(t *testing.T) {
	testPanicLog()
	Info(context.Background(), "catch panic")
}

func testPanicLog() {
	ctx := context.WithValue(context.Background(), LOG_TRACE_ID, "zap_trace_id")

	defer Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("recover from err: %v", e)
	})

	panic(errors.New("dddd"))
}

func TestLevelChange(t *testing.T) {
	if LogLess() != zap.WarnLevel {
		t.FailNow()
	}

	if LogLess() != zap.ErrorLevel {
		t.FailNow()
	}

	if LogLess() != zap.DPanicLevel {
		t.FailNow()
	}

	if LogLess() != zap.PanicLevel {
		t.FailNow()
	}

	if LogLess() != zap.FatalLevel {
		t.FailNow()
	}

	if LogLess() != zap.FatalLevel {
		t.FailNow()
	}

	if LogMore() != zap.PanicLevel {
		t.FailNow()
	}

	if LogMore() != zap.DPanicLevel {
		t.FailNow()
	}

	if LogMore() != zap.ErrorLevel {
		t.FailNow()
	}

	if LogMore() != zap.WarnLevel {
		t.FailNow()
	}

	if LogMore() != zap.InfoLevel {
		t.FailNow()
	}

	if LogMore() != zap.DebugLevel {
		t.FailNow()
	}

	if LogMore() != zap.DebugLevel {
		t.FailNow()
	}
}
