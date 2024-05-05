**Read this in other languages: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [log](#log)
  * [zap_test.go](#zap_testgo)

<!-- tocstop -->

# log
## zap_test.go
### TestZapConsole
```go

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
```
### TestZapJson
```go

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
```
### TestZapHttpOut
```go

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
```
### TestLevelChange
```go

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
```
### TestNewTraceId
```go

t.Log(NewTraceId())
```
