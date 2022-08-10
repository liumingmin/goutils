

<!-- toc -->

- [log zap日志库](#log-zap%E6%97%A5%E5%BF%97%E5%BA%93)
  * [zap_test.go](#zap_testgo)
    + [TestZap](#testzap)
    + [TestPanicLog](#testpaniclog)
    + [TestLevelChange](#testlevelchange)

<!-- tocstop -->

# log zap日志库
## zap_test.go
### TestZap
```go

ctx := &gin.Context{}
ctx.Set(LOG_TRADE_ID, "aaabbbbbcccc")

Info(ctx, "我是日志2")
Error(ctx, "我是日志4: %v,%v", "管理员", "eee")
```
### TestPanicLog
```go

testPanicLog()

Info(context.Background(), "catch panic")
```
### TestLevelChange
```go

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
```
