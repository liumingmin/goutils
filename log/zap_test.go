package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GameDefaultFieldGenerator struct {
}

func (f *GameDefaultFieldGenerator) Generate(ctx context.Context) []zapcore.Field {
	return []zap.Field{
		zap.String("gameCode", "lol"),
		zap.String("version", "1.0"),
	}
}

func TestZap(t *testing.T) {
	testRunLogServer(t, "log1", "log2")

	BaseFieldsGenerator = &GameDefaultFieldGenerator{}
	SetLogLevel(zap.DebugLevel)

	ctx := ContextWithTraceId()
	Debug(ctx, "I am debug log1")
	LogLess()
	Debug(ctx, "I am debug log2")

	Info(ctx, "I am info log1: %v", "hello")
	LogLess()
	Info(ctx, "I am info log2: %v", "hello")

	Warn(ctx, "I am warn log1: %v", "hello sam")
	LogLess()
	Warn(ctx, "I am warn log2: %v", "hello sam")

	Error(ctx, "I am error log1: %v, %v", "admin", "eee")
	LogLess()
	Error(ctx, "I am error log2: %v, %v", "admin", "eee")

	ErrorStack(ctx, "I am panic error log1")
	LogLess()
	ErrorStack(ctx, "I am panic error log2")

	Log(ctx, zap.DebugLevel, "I am debug, log1")

	testPanicLog(func() {
		panic(errors.New("I am log1"))
	})

	time.Sleep(time.Second)
}

func testRunLogServer(t *testing.T, incldueTag, excludeTag string) {
	http.HandleFunc("/goutils/log", func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		str := string(data)

		if !strings.Contains(str, "I am") {
			return
		}

		if incldueTag != "" && !strings.Contains(str, incldueTag) {
			t.FailNow()
		}

		if excludeTag != "" && strings.Contains(str, excludeTag) {
			t.FailNow()
		}
	})
	go http.ListenAndServe(":8053", nil)
}

func testPanicLog(f func()) {
	ctx := ContextWithTraceId()

	defer Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("catch panic: %v", e)
	})

	f()
}

func TestLevelChange(t *testing.T) {
	SetLogLevel(zap.DebugLevel)

	if GetLogLevel() != zap.DebugLevel {
		t.FailNow()
	}

	if LogLess() != zap.InfoLevel {
		t.FailNow()
	}

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

	SetLogLevel(zap.FatalLevel)

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

func TestNewTraceId(t *testing.T) {
	t.Log(NewTraceId())
}
