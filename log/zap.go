package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/liumingmin/goutils/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LOG_TRACE_ID            = "__GTraceId__"
	LOG_JSON_FIELD_TRACE_ID = "traceId"

	LOGGER_ENCODER_JSON    = "json"
	LOGGER_ENCODER_CONSOLE = "console"
)

type ZapLogImpl struct {
	logger      *zap.Logger
	stackLogger *zap.Logger

	loggerLevel zap.AtomicLevel
}

func NewZapLogImpl() ILog {
	loggerLevel := populateLogLevel()

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "log",
		CallerKey:     "linenum",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}, // 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	var cores []zapcore.Core
	for _, logConf := range conf.Conf.Logs {
		cores = append(cores, zapcore.NewCore(
			populateEncoder(logConf)(encoderConfig),                      // 编码器配置
			zapcore.NewMultiWriteSyncer(populateWriteSyncer(logConf)...), // 打印到控制台、文件、HTTP
			loggerLevel, // 日志级别
		))
	}

	// 构造日志
	logger := zap.New(zapcore.NewTee(cores...), caller, development, zap.AddCallerSkip(1))
	stackLogger := logger.WithOptions(zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))

	return &ZapLogImpl{
		logger:      logger,
		stackLogger: stackLogger,
		loggerLevel: loggerLevel,
	}
}

func populateEncoder(logConf *conf.Log) func(zapcore.EncoderConfig) zapcore.Encoder {
	// JSON 编码器
	if logConf.OutputEncoder == LOGGER_ENCODER_JSON {
		return zapcore.NewJSONEncoder
	}
	// 默认：CONSOLE 编码器
	return zapcore.NewConsoleEncoder
}

func populateWriteSyncer(logConf *conf.Log) []zapcore.WriteSyncer {
	writeSyncers := make([]zapcore.WriteSyncer, 0)
	// 标准输出流
	if logConf.Stdout {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}
	// 文件输出流
	if logConf.FileOut {
		populateLogHook(logConf)
		writeSyncers = append(writeSyncers, zapcore.AddSync(&logConf.Logger))
	}
	// Http输出流
	if logConf.HttpOut {
		hWriter := &httpWriter{
			logConf:          logConf,
			loggerHttpClient: &http.Client{Timeout: time.Second * time.Duration(logConf.HttpTimeout)},
		}
		writeSyncers = append(writeSyncers, zapcore.AddSync(hWriter))
	}
	return writeSyncers
}

func populateLogHook(logConf *conf.Log) {
	if logConf.Logger.Filename == "" {
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
		logConf.Logger.Filename = logFileName // 日志文件路径
	}
	if logConf.Logger.MaxSize == 0 {
		logConf.Logger.MaxSize = 128 // 每个日志文件保存的最大尺寸 单位：M
	}
	if logConf.Logger.MaxBackups == 0 {
		logConf.Logger.MaxBackups = 7
	}
	if logConf.Logger.MaxAge == 0 {
		logConf.Logger.MaxAge = 7
	}
}

func populateLogLevel() zap.AtomicLevel {
	loggerLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	if conf.Conf.LogLevel != "" {
		loggerLevel.UnmarshalText([]byte(conf.Conf.LogLevel))
	}
	return loggerLevel
}

func CnTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func (l *ZapLogImpl) Log(ctx context.Context, level zapcore.Level, msg string, args ...interface{}) {
	if !l.logger.Core().Enabled(level) {
		return
	}

	if len(args) > 0 {
		switch args[0].(type) {
		case zapcore.Field:
			l.logger.Log(level, msg, l.argsToZapFields(ctx, args...)...)
		default:
			l.logger.Log(level, l.argsToMsg(ctx, msg, args...), BaseFieldsGenerator.Generate(ctx)...)
		}
		return
	}
	l.logger.Log(level, msg)
}

func (l *ZapLogImpl) Debug(ctx context.Context, msg string, args ...interface{}) {
	if !l.logger.Core().Enabled(zap.DebugLevel) {
		return
	}

	if len(args) > 0 {
		switch args[0].(type) {
		case zapcore.Field:
			l.logger.Debug(msg, l.argsToZapFields(ctx, args...)...)
		default:
			l.logger.Debug(l.argsToMsg(ctx, msg, args...), BaseFieldsGenerator.Generate(ctx)...)
		}
		return
	}

	l.logger.Debug(msg)
}

func (l *ZapLogImpl) Info(ctx context.Context, msg string, args ...interface{}) {
	if !l.logger.Core().Enabled(zap.InfoLevel) {
		return
	}

	if len(args) > 0 {
		switch args[0].(type) {
		case zapcore.Field:
			l.logger.Info(msg, l.argsToZapFields(ctx, args...)...)
		default:
			l.logger.Info(l.argsToMsg(ctx, msg, args...), BaseFieldsGenerator.Generate(ctx)...)
		}
		return
	}

	l.logger.Info(msg)
}

func (l *ZapLogImpl) Warn(ctx context.Context, msg string, args ...interface{}) {
	if !l.logger.Core().Enabled(zap.WarnLevel) {
		return
	}

	if len(args) > 0 {
		switch args[0].(type) {
		case zapcore.Field:
			l.logger.Warn(msg, l.argsToZapFields(ctx, args...)...)
		default:
			l.logger.Warn(l.argsToMsg(ctx, msg, args...), BaseFieldsGenerator.Generate(ctx)...)
		}
		return
	}

	l.logger.Warn(msg)
}

func (l *ZapLogImpl) Error(ctx context.Context, msg string, args ...interface{}) {
	if !l.logger.Core().Enabled(zap.ErrorLevel) {
		return
	}

	if len(args) > 0 {
		switch args[0].(type) {
		case zapcore.Field:
			l.logger.Error(msg, l.argsToZapFields(ctx, args...)...)
		default:
			l.logger.Error(l.argsToMsg(ctx, msg, args...), BaseFieldsGenerator.Generate(ctx)...)
		}
		return
	}

	l.logger.Error(msg)
}

func (l *ZapLogImpl) Recover(ctx context.Context, errHandler func(interface{}) string) {
	if err := recover(); err != nil {
		traceStr := ""
		traceId := ctx.Value(LOG_TRACE_ID)
		if traceId != nil {
			traceStr = "<" + fmt.Sprint(traceId) + "> "
		}

		l.stackLogger.Error(traceStr+" panic: "+errHandler(err), BaseFieldsGenerator.Generate(ctx)...)
	}
}

func (l *ZapLogImpl) ErrorStack(ctx context.Context, msg string, args ...interface{}) {
	if !l.stackLogger.Core().Enabled(zap.ErrorLevel) {
		return
	}

	if len(args) > 0 {
		switch args[0].(type) {
		case zapcore.Field:
			l.stackLogger.Error(msg, l.argsToZapFields(ctx, args...)...)
		default:
			l.stackLogger.Error(l.argsToMsg(ctx, msg, args...), BaseFieldsGenerator.Generate(ctx)...)
		}
		return
	}

	l.stackLogger.Error(msg)
}

func (l *ZapLogImpl) argsToMsg(ctx context.Context, msg string, args ...interface{}) string {
	traceStr := ""
	traceId := ctx.Value(LOG_TRACE_ID)
	if traceId != nil {
		traceStr = "<" + fmt.Sprint(traceId) + "> "
	}

	return traceStr + fmt.Sprintf(msg, args...)
}

func (l *ZapLogImpl) argsToZapFields(ctx context.Context, args ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(args)+1)
	for _, arg := range args {
		fields = append(fields, arg.(zap.Field))
	}

	traceId := ctx.Value(LOG_TRACE_ID)
	if traceId != nil {
		fields = append(fields, zap.String(LOG_JSON_FIELD_TRACE_ID, fmt.Sprint(traceId)))
	}

	if BaseFieldsGenerator != nil {
		fs := BaseFieldsGenerator.Generate(ctx)
		if fs != nil {
			fields = append(fields, fs...)
		}
	}
	return fields
}

func (l *ZapLogImpl) GetLogLevel() zapcore.Level {
	return l.loggerLevel.Level()
}

func (l *ZapLogImpl) SetLogLevel(lvl zapcore.Level) zapcore.Level {
	oldLevel := l.loggerLevel.Level()
	l.loggerLevel.SetLevel(lvl)
	return oldLevel
}

func (l *ZapLogImpl) LogMore() zapcore.Level {
	level := l.loggerLevel.Level()
	if level == zap.DebugLevel {
		return level
	}
	l.loggerLevel.SetLevel(level - 1)
	return level - 1
}

func (l *ZapLogImpl) LogLess() zapcore.Level {
	level := l.loggerLevel.Level()
	if level == zap.FatalLevel {
		return level
	}
	l.loggerLevel.SetLevel(level + 1)
	return level + 1
}

type DefaultFieldsGenerator struct {
	Nop []zapcore.Field
}

func (f *DefaultFieldsGenerator) Generate(ctx context.Context) []zapcore.Field {
	return f.Nop
}

type httpWriter struct {
	logConf          *conf.Log
	loggerHttpClient *http.Client
}

func (h *httpWriter) Write(data []byte) (int, error) {
	if h.logConf.HttpUrl == "" {
		return 0, nil
	}

	input := make([]byte, len(data))
	copy(input, data)

	go func() {
		defer func() {
			if e := recover(); e != nil {
				if h.logConf.HttpDebug {
					fmt.Fprintf(os.Stderr, "http log failed, err: %v", e)
				}
			}
		}()

		resp, err := h.loggerHttpClient.Post(h.logConf.HttpUrl, "application/json", bytes.NewBuffer(input))
		if err != nil {
			if h.logConf.HttpDebug {
				fmt.Fprintf(os.Stderr, "http log failed, err: %+v, data: %+v", err, string(input))
			}
			return
		}
		defer resp.Body.Close()

		if h.logConf.HttpDebug {
			body, _ := io.ReadAll(resp.Body)
			fmt.Fprintf(os.Stdout, "http log successful: %+v, data: %+v", string(body), string(input))
		}
	}()

	return len(input), nil
}
