package log

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/liumingmin/goutils/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type GameDefaultFieldGenerator struct {
}

func (f *GameDefaultFieldGenerator) Generate(ctx context.Context) []zapcore.Field {
	return []zap.Field{
		zap.String("gameCode", "lol"),
		zap.String("version", "1.0"),
	}
}

func (f *GameDefaultFieldGenerator) GenerateStr(ctx context.Context) string {
	return fmt.Sprintf("gameCode: %v, version: %v", "lol", "1.0")
}

func TestZapConsole(t *testing.T) {
	conf.Conf.LogLevel = "debug"
	conf.Conf.Logs = []*conf.Log{{
		Logger:        lumberjack.Logger{Filename: "goutils.log"},
		OutputEncoder: "console",
		Stdout:        false,
		FileOut:       true,
	}}
	LogImpl = NewZapLogImpl()

	ctx := ContextWithTraceId()
	Debug(ctx, "I am debug log1")
	LogLess()
	Debug(ctx, "I am debug log2")

	Info(ctx, "I am info log1: %v", "hello")
	Log(ctx, zap.InfoLevel, "I am info level log1")
	LogLess()
	Log(ctx, zap.InfoLevel, "I am info level log2")
	Info(ctx, "I am info log2: %v", "hello")

	Warn(ctx, "I am warn log1: %v", "hello tom")
	Log(ctx, zap.WarnLevel, "I am warn level log1: %v", "hello sam")
	LogLess()
	Warn(ctx, "I am warn log2: %v", "hello tom")
	Log(ctx, zap.WarnLevel, "I am warn level log2: %v", "hello sam")

	testPanicLog(func() {
		time.Sleep(time.Millisecond * 10)
		// panic(errors.New("test panic"))
	})

	Error(ctx, "I am error log1: %v, %v", "admin", "eee")
	ErrorStack(ctx, "test panic log1")

	LogLess()

	Error(ctx, "I am error log2: %v, %v", "admin", "eee")
	ErrorStack(ctx, "test panic log2")
}

func testPanicLog(f func()) {
	ctx := ContextWithTraceId()

	defer Recover(ctx, func(e interface{}) string {
		return fmt.Sprintf("catch panic: %v", e)
	})

	f()
}

func TestZapJson(t *testing.T) {
	conf.Conf.Logs = []*conf.Log{{
		Logger:        lumberjack.Logger{Filename: "goutils.json"},
		OutputEncoder: "json",
		Stdout:        true,
		FileOut:       true,
	}}
	LogImpl = NewZapLogImpl()

	ctx := ContextWithTraceId()

	BaseFieldsGenerator = &GameDefaultFieldGenerator{}
	Debug(ctx, "I am debug log1")

	Info(ctx, "I am info log1", zap.String("userId", "0001"), zap.Bool("isAdmin", true))

	Log(ctx, zap.WarnLevel, "I am warn log1", zap.String("userId", "0002"), zap.Bool("isAdmin", false))
}

func TestZapHttpOut(t *testing.T) {
	conf.Conf.LogLevel = "debug"
	conf.Conf.Logs = []*conf.Log{{
		OutputEncoder: "console",
		Stdout:        true,
		FileOut:       true,
		HttpOut:       true,
		HttpUrl:       "http://127.0.0.1:8053/goutils/log",
		HttpDebug:     true,
	}}
	LogImpl = NewZapLogImpl()
	testRunLogServer(t, "log1", "log2")

	ctx := ContextWithTraceId()
	Debug(ctx, "I am debug log1")
	LogLess()
	Debug(ctx, "I am debug log2")

	time.Sleep(time.Second)
}

func testRunLogServer(t *testing.T, incldueTag, excludeTag string) {
	http.HandleFunc("/goutils/log", func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		str := string(data)
		//fmt.Println(str)
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
