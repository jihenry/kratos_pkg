package zaplog

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	rlog "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLoggerConf struct {
	Level       string `yaml:"level"` //日志等级 debug->info->warn->error
	FileName    string `yaml:"fileName"`
	FilePath    string `yaml:"filePath"`
	MaxSize     int64  `yaml:"maxSize"`     //单位为mb
	MaxBackups  int64  `yaml:"maxBackups"`  //最多保留日志备份
	MaxAge      int64  `yaml:"maxAge"`      //备份最大生命周期 0为长期保存 单位：天
	Compress    bool   `yaml:"compress"`    //是否压缩
	ShowConsole bool   `yaml:"showConsole"` //同时输出到控制台
}

func writer(zc ZapLoggerConf) (io.Writer, error) {
	month := time.Now().Format("200601")
	filePath := filepath.Join(zc.FilePath, month)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, 0755)
	}
	return rlog.New(
		filePath+"/"+zc.FileName+".%Y%m%d.log",
		rlog.WithLinkName(zc.FilePath+zc.FileName+".log"),
		rlog.WithMaxAge(time.Hour*24*time.Duration(zc.MaxAge)),
		rlog.WithRotationTime(time.Hour*24),
	)
}

func newZapLogger(zc ZapLoggerConf) *zap.SugaredLogger {
	writer, err := writer(zc)
	if err != nil {
		panic(err)
	}
	//是否输出到控制台
	var (
		syncer zapcore.WriteSyncer
	)
	if zc.ShowConsole {
		syncer = zapcore.NewMultiWriteSyncer(
			//标准输出
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(writer),
		)
	} else {
		syncer = zapcore.AddSync(writer)
	}

	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
	}
	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// 自定义文件：行号输出项
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		//enc.AppendString("[" + l.traceId + "]") // 链路追踪id
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "line",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      customLevelEncoder, // 小写编码器
		EncodeTime:       customTimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder, //
		EncodeCaller:     customCallerEncoder,            // 全路径编码器
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " | ",
	}

	var (
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
		level   = getLevel(zc.Level)
	)
	core := zapcore.NewCore(
		encoder,
		syncer,
		level,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	return logger

}

func getLevel(l string) zapcore.Level {
	level := zap.InfoLevel
	switch l {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}
	return level
}

//const (
//	loggerCtxKey = "Ctx-Key-Logger"
//)
//
//func NewLoggerContext() string {
//	return loggerCtxKey
//}

type LoggerCtxKey string

func Loggers() *zap.SugaredLogger {
	if sugarLogger == nil {
		panic("zapLogger is nil")
	}
	return sugarLogger
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		panic("nil ctx")
	}
	return FromLoggerContext(ctx)
}

type loggerKey struct {
}

func NewLoggerContext(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromLoggerContext(ctx context.Context) *zap.SugaredLogger {
	if logger, ok := ctx.Value(loggerKey{}).(*zap.SugaredLogger); ok {
		return logger
	}
	return Loggers()
}

// Level is a logger level.
type Level int8

// LevelKey is logger level key.
const LevelKey = "level"

const (
	// LevelDebug is logger debug level.
	LevelDebug Level = iota - 1
	// LevelInfo is logger info level.
	LevelInfo
	// LevelWarn is logger warn level.
	LevelWarn
	// LevelError is logger error level.
	LevelError
	// LevelFatal is logger fatal level
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// ParseLevel parses a level string into a logger Level value.
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	}
	return LevelInfo
}
