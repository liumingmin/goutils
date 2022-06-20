package log

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/liumingmin/goutils/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LOG_TRADE_ID = "__GTraceId__"

var (
	logger      *zap.Logger
	loggerLevel zap.AtomicLevel
	stackLogger *zap.Logger
)

func init() {
	hook := conf.Conf.Log.Logger

	if hook.Filename == "" {
		file, _ := exec.LookPath(os.Args[0])
		filename := filepath.Base(file)
		extName := filepath.Ext(filename)
		logFileName := ""
		if extName != "" {
			logFileName = strings.Replace(filename, extName, ".log", -1)
		} else {
			logFileName = filename + ".log"
		}
		hook.Filename = logFileName // 日志文件路径
	}
	if hook.MaxSize == 0 {
		hook.MaxSize = 128 // 每个日志文件保存的最大尺寸 单位：M
	}
	if hook.MaxBackups == 0 {
		hook.MaxBackups = 7
	}
	if hook.MaxAge == 0 {
		hook.MaxAge = 7
	}

	// 设置日志级别
	loggerLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	if conf.Conf.Log.LogLevel != "" {
		loggerLevel.UnmarshalText([]byte(conf.Conf.Log.LogLevel))
	}

	writeSyncers := make([]zapcore.WriteSyncer, 0)
	if conf.Conf.Log.Stdout {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}
	if conf.Conf.Log.FileOut {
		writeSyncers = append(writeSyncers, zapcore.AddSync(&hook))
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "log",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 小写编码器
		EncodeTime:     CnTimeEncoder,                  // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),     // 编码器配置 NewConsoleEncoder NewJSONEncoder
		zapcore.NewMultiWriteSyncer(writeSyncers...), // 打印到控制台和文件
		loggerLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	logger = zap.New(core, caller, development, zap.AddCallerSkip(1))
	stackLogger = logger.WithOptions(zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))

	Debug(context.Background(), "log 初始化成功")
}

func CnTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func Log(c context.Context, level zapcore.Level, args ...interface{}) {
	if !logger.Core().Enabled(level) {
		return
	}

	msg := parseArgs(c, args...)
	if ce := logger.Check(level, msg); ce != nil {
		ce.Write()
	}
}

func Debug(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.DebugLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Debug(msg)
}

func Info(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.InfoLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Info(msg)
}

func Warn(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.WarnLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Warn(msg)
}

func Error(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.ErrorLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Error(msg)
}

func Fatal(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.FatalLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Fatal(msg)
}

func Panic(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.PanicLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Panic(msg)
}

func LogMore() zapcore.Level {
	level := loggerLevel.Level()
	if level == zap.DebugLevel {
		return level
	}
	loggerLevel.SetLevel(level - 1)
	return level - 1
}

func LogLess() zapcore.Level {
	level := loggerLevel.Level()
	if level == zap.FatalLevel {
		return level
	}
	loggerLevel.SetLevel(level + 1)
	return level + 1
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
		var ok bool
		msg, ok = args[0].(string)
		if !ok {
			msg = fmt.Sprint(args[0])
		}

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
	traceId := c.Value(LOG_TRADE_ID)
	if traceId != nil {
		return "<" + fmt.Sprint(traceId) + ">"
	}

	return ""
}
