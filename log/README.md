**Read this in other languages: [English](README.md), [中文](README_zh.md).**



<!-- toc -->

- [log](#log)
  * [zap_test.go](#zap_testgo)

<!-- tocstop -->

# log
## zap_test.go
### TestZap
```go

ctx := ContextWithTraceId()
Info(ctx, "我是日志2")
SetDefaultGenerator(new(GameDefaultFieldGenerator))
Error(ctx, "我是日志4: %v,%v", "管理员", "eee")
Info(ctx, "我是日志5: %v", "hello")
Warn(ctx, "我是日志6: %v", "hello sam")
```
### TestErrorStack
```go

ErrorStack(context.Background(), "panic error")
```
### TestPanicLog
```go

testPanicLog()
Info(context.Background(), "catch panic")
```
### TestLevelChange
```go

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
```
### TestNewTraceId
```go

t.Log(NewTraceId())
```
