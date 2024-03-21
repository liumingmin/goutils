package log

import (
	"bytes"
	"context"
	"encoding/base32"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/liumingmin/goutils/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	LOG_TRACE_ID            = "__GTraceId__"
	LOG_JSON_FIELD_TRACE_ID = "traceId"

	LOGGER_ENCODER_JSON    = "json"
	LOGGER_ENCODER_CONSOLE = "console"
)

var (
	logger           *zap.Logger
	loggerLevel      zap.AtomicLevel
	stackLogger      *zap.Logger
	loggerHttpClient *http.Client
	generator        DefaultFieldsGenerator
	generatorLock    sync.Mutex
)

func init() {
	generator = new(DefaultGenerator)

	syncers := populateWriteSyncer()

	// 创建zap core
	core := zapcore.NewCore(
		populateEncoder(),                       // 编码器配置 NewConsoleEncoder NewJSONEncoder
		zapcore.NewMultiWriteSyncer(syncers...), // 打印到控制台、文件、HTTP
		populateLogLevel(),                      // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	logger = zap.New(core, caller, development, zap.AddCallerSkip(1))
	stackLogger = logger.WithOptions(zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))

	loggerHttpClient = &http.Client{Timeout: time.Second * time.Duration(conf.Conf.Log.HttpTimeout)}

	Debug(context.Background(), "log 初始化成功")
}

func populateWriteSyncer() []zapcore.WriteSyncer {
	writeSyncers := make([]zapcore.WriteSyncer, 0)
	// 标准输出流
	if conf.Conf.Log.Stdout {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}
	// 文件输出流
	if conf.Conf.Log.FileOut {
		hook := populateLogHook()
		writeSyncers = append(writeSyncers, zapcore.AddSync(&hook))
	}
	// Http输出流
	if conf.Conf.Log.HttpOut {
		writeSyncers = append(writeSyncers, zapcore.AddSync(new(httpWriter)))
	}
	return writeSyncers
}

func populateEncoder() zapcore.Encoder {
	config := zapcore.EncoderConfig{
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
	// JSON 编码器
	if conf.Conf.Log.OutputEncoder == LOGGER_ENCODER_JSON {
		return zapcore.NewJSONEncoder(config)
	}
	// 默认：CONSOLE 编码器
	return zapcore.NewConsoleEncoder(config)
}

func populateLogHook() lumberjack.Logger {
	hook := conf.Conf.Log.Logger

	if hook.Filename == "" {
		file, _ := exec.LookPath(os.Args[0])
		filename := filepath.Base(file)
		extName := filepath.Ext(filename)
		logFileName := ""
		if extName != "" {
			extIndex := strings.LastIndex(filename, extName)
			logFileName = filename[:extIndex] + ".log"
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
	return hook
}

func populateLogLevel() zapcore.LevelEnabler {
	loggerLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	if conf.Conf.Log.LogLevel != "" {
		loggerLevel.UnmarshalText([]byte(conf.Conf.Log.LogLevel))
	}
	return loggerLevel
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
	logger.Debug(msg, SupplementFields(c)...)
}

func Info(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.InfoLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Info(msg, SupplementFields(c)...)
}

func Warn(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.WarnLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Warn(msg, SupplementFields(c)...)
}

func Error(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.ErrorLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Error(msg, SupplementFields(c)...)
}

func Fatal(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.FatalLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Fatal(msg, SupplementFields(c)...)
}

func Panic(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.PanicLevel) {
		return
	}

	msg := parseArgs(c, args...)
	logger.Panic(msg, SupplementFields(c)...)
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

func Recover(c context.Context, errHandler func(interface{}) string) {
	if err := recover(); err != nil {
		stackLogger.Error(ctxParams(c)+" panic: "+errHandler(err), SupplementFields(c)...)
	}
}

func ErrorStack(c context.Context, args ...interface{}) {
	if !logger.Core().Enabled(zap.ErrorLevel) {
		return
	}

	msg := parseArgs(c, args...)
	stackLogger.Error(msg)
}

func parseArgs(c context.Context, args ...interface{}) (msg string) {
	var paramArgs []interface{}
	if len(args) == 0 {
		msg = ""
	} else {
		var ok bool
		msg, ok = args[0].(string)
		if !ok {
			msg = fmt.Sprint(args[0])
		}

		if len(args) > 1 {
			paramArgs = args[1:]
		}
	}

	if len(paramArgs) > 0 {
		msg = fmt.Sprintf(msg, paramArgs...)
	}

	msg = ctxParams(c) + msg
	return msg
}

func ctxParams(c context.Context) string {
	if conf.Conf.Log.OutputEncoder == LOGGER_ENCODER_JSON {
		return ""
	}

	traceId := c.Value(LOG_TRACE_ID)
	if traceId != nil {
		return "<" + fmt.Sprint(traceId) + "> "
	}

	return ""
}

func SupplementFields(ctx context.Context) []zap.Field {
	if conf.Conf.Log.OutputEncoder != LOGGER_ENCODER_JSON {
		return generator.GetDefaultFields()
	}

	traceId := ctx.Value(LOG_TRACE_ID)
	if traceId != nil {
		return append(generator.GetDefaultFields(), zap.String(LOG_JSON_FIELD_TRACE_ID, fmt.Sprint(traceId)))
	}
	return generator.GetDefaultFields()
}

// DefaultFieldsGenerator 默认值入参
type DefaultFieldsGenerator interface {
	GetDefaultFields() []zap.Field
}

type DefaultGenerator struct {
}

func (f *DefaultGenerator) GetDefaultFields() []zap.Field {
	return nil
}

func SetDefaultGenerator(g DefaultFieldsGenerator) {
	generatorLock.Lock()
	defer generatorLock.Unlock()
	generator = g
}

type httpWriter struct {
}

func (h *httpWriter) Write(data []byte) (int, error) {
	if conf.Conf.Log.HttpUrl == "" {
		return 0, nil
	}

	input := make([]byte, len(data))
	copy(input, data)

	go func() {
		defer func() {
			if e := recover(); e != nil {
				if conf.Conf.Log.HttpDebug {
					fmt.Fprintf(os.Stderr, "http log failed, err: %v", e)
				}
			}
		}()

		resp, err := loggerHttpClient.Post(conf.Conf.Log.HttpUrl, "application/json", bytes.NewBuffer(input))
		if err != nil {
			if conf.Conf.Log.HttpDebug {
				fmt.Fprintf(os.Stderr, "http log failed, err: %+v, data: %+v", err, string(input))
			}
			return
		}
		defer resp.Body.Close()

		if conf.Conf.Log.HttpDebug {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Fprintf(os.Stdout, "http log successful: %+v, data: %+v", string(body), string(input))
		}
	}()

	return len(input), nil
}

func ContextWithTraceId() context.Context {
	return ContextWithTraceIdFromParent(context.Background())
}

func ContextWithTraceIdFromParent(parent context.Context) context.Context {
	uuidBytes := [16]byte(uuid.New())
	b32Uuid := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(uuidBytes[:])
	return context.WithValue(parent, LOG_TRACE_ID, b32Uuid)
}
