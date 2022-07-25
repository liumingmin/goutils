# log zap日志库
## zap_test.go
### TestZap
```go

ctx := &gin.Context{}
ctx.Set("__traceId", "aaabbbbbcccc")
//Info(ctx, "我是日志", "name", "管理员")  //json

Info(ctx, "我是日志2")

//Info(ctx, "我是日志3", "name")  //json

Error(ctx, "我是日志4: %v,%v", "管理员", "eee")
```
### TestZapJson
```go

ctx := &gin.Context{}
ctx.Set("__traceId", "aaabbbbbcccc")
Info(ctx, "我是日志 %v", "name", "管理员") //json

Info(ctx, "我是日志3 %v", "管理员") //json
Error(ctx, "我是日志3")          //json
Log(ctx, zapcore.ErrorLevel, "日志啊")
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
