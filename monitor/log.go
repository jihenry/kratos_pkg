package monitor

import (
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogConfig zaplog配置
type ZapLogConfig struct {
	FileName    string `yaml:"FileName"`
	FilePath    string `yaml:"FilePath"`
	MaxSize     int64  `yaml:"MaxSize"`     //单位为mb
	MaxBackups  int64  `yaml:"MaxBackups"`  //最多保留日志备份
	MaxAge      int64  `yaml:"MaxAge"`      //备份最大生命周期 0为长期保存 单位：天
	Compress    bool   `yaml:"Compress"`    //是否压缩
	ShowConsole bool   `yaml:"ShowConsole"` //同时输出到控制台
}

var log *zap.SugaredLogger

func initLog(fileName, logFilePath string) {
	c := ZapLogConfig{FileName: fileName, FilePath: logFilePath, MaxSize: 100, MaxBackups: 30, MaxAge: 30, Compress: true, ShowConsole: true}
	log = NewLogger(c)
}

// InitLogger 初始化配置
func InitLogger(confLog ZapLogConfig) {
	log = NewLogger(confLog)
}

// Loggers 获得日志
func Loggers() *zap.SugaredLogger {
	if log == nil {
		panic("nil sugarLogger")
	}
	return log
}

// NewLogger 声明一个新配置
func NewLogger(confLog ZapLogConfig) *zap.SugaredLogger {
	writer := getWriter(confLog)

	var syncer zapcore.WriteSyncer
	if confLog.ShowConsole { //判断是否要打印到控制台
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(writer))
	} else {
		syncer = zapcore.AddSync(writer)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    nil,
		EncodeTime:     nil,
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   nil,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder = zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(
		encoder,
		syncer,
		zap.InfoLevel,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	return logger
}

func getWriter(confLog ZapLogConfig) io.Writer {
	_ = os.MkdirAll(confLog.FilePath, 0755)

	hook, err := rotatelogs.New(
		filepath.Join(confLog.FilePath, confLog.FileName+".%Y%m%d"+".log"),
		rotatelogs.WithLinkName(filepath.Join(confLog.FilePath, confLog.FileName+".log")),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(confLog.MaxAge)),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}
