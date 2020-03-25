package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"goutils/conf"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger      *zap.Logger
	stackLogger *zap.Logger
)

const (
	GlobalTraceId = "__gTraceId"
)

func init() {
	hook := lumberjack.Logger{
		Filename:   "./goutils.log", // 日志文件路径
		MaxSize:    128,             // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,              // 日志文件最多保存多少个备份
		MaxAge:     7,               // 文件最多保存多少天
		Compress:   true,            // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "log",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 小写编码器
		EncodeTime:     CnTimeEncoder,                    // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,   //
		EncodeCaller:   zapcore.FullCallerEncoder,        // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()

	if debug := conf.ExtBool("debug", true); debug {
		atomicLevel.SetLevel(zap.InfoLevel)
	} else {
		atomicLevel.SetLevel(zap.ErrorLevel)
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                                        // 编码器配置 NewConsoleEncoder NewJSONEncoder
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	logger = zap.New(core, caller, development, zap.AddCallerSkip(1))
	stackLogger = logger.WithOptions(zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))

	Info(context.Background(), "log 初始化成功")
}

func CnTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func Debug(c context.Context, args ...interface{}) {
	msg := parseArgs(c, args...)
	logger.Debug(msg)
}

func Info(c context.Context, args ...interface{}) {
	msg := parseArgs(c, args...)
	logger.Info(msg)
}

func Warn(c context.Context, args ...interface{}) {
	msg := parseArgs(c, args...)
	logger.Warn(msg)
}

func Error(c context.Context, args ...interface{}) {
	msg := parseArgs(c, args...)
	logger.Error(msg)
}

func Fatal(c context.Context, args ...interface{}) {
	msg := parseArgs(c, args...)
	logger.Fatal(msg)
}

func Panic(c context.Context, args ...interface{}) {
	msg := parseArgs(c, args...)
	logger.Panic(msg)
}

func Recover(c context.Context, arg0 interface{}) {
	recoverArgs := []interface{}{"%v %v"}
	if err := recover(); err != nil {
		switch first := arg0.(type) {
		case func(interface{}) string:
			recoverArgs = append(recoverArgs, []interface{}{"error", first(err)}...)

		default:
			recoverArgs = append(recoverArgs, []interface{}{"error", err}...)
		}

		msg := parseArgs(c, recoverArgs...)
		stackLogger.Error(msg)
	}
}

func parseArgs(c context.Context, args ...interface{}) (msg string) {
	parmArgs := make([]interface{}, 0)
	if len(args) == 0 {
		msg = ""
	} else {
		msg = fmt.Sprint(args[0])
		parmArgs = args[1:]
	}

	lenParmArgs := len(parmArgs)

	if lenParmArgs > 0 {
		msg = fmt.Sprintf(msg, parmArgs...)
	}

	msg = ctxParams(c) + " " + msg

	return
}

func ctxParams(c context.Context) string {
	traceId := c.Value("__traceId")
	if traceId != nil {
		return "<" + fmt.Sprint(traceId) + ">"
	}

	return ""
}
