**Read this in other languages: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [log](#log)
  * [zap_test.go](#zap_testgo)

<!-- tocstop -->

# log
## zap_test.go
### TestZap
```go

testRunLogServer(t, "log1", "log2")
SetDefaultGenerator(new(GameDefaultFieldGenerator))
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

testPanicLog(func() {
	Panic(context.Background(), "I am panic log1")
})
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
