package log

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*options)

type options struct {
	dir      string                //存放目录
	fileName string                //文件名
	console  bool                  //是否显示在终端，用于命令行启动查看日志
	level    string                //用于控制什么等级的日志才输出
	encoder  zapcore.EncoderConfig //用于外部扩展
	zapOpts  []zap.Option          //用于外部扩展选项
	maxAge   int                   //保存最大天数, 0为长期保存
}

func WithDir(dir string) Option {
	return func(o *options) {
		o.dir = dir
	}
}

func WithFileName(fname string) Option {
	return func(o *options) {
		o.fileName = fname
	}
}

func WithConsole(console bool) Option {
	return func(o *options) {
		o.console = console
	}
}

func WithEncoder(encoder zapcore.EncoderConfig) Option {
	return func(o *options) {
		o.encoder = encoder
	}
}

func WithZapOpts(opts ...zap.Option) Option {
	return func(o *options) {
		o.zapOpts = opts
	}
}

func WithLevel(level string) Option {
	return func(o *options) {
		o.level = level
	}
}

var _ log.Logger = (*zapLogWraper)(nil)

type zapLogWraper struct {
	log *zap.Logger
}

func NewZapLogger(opts ...Option) (log.Logger, error) {
	defaultEncoderCfg := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stack",
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: zapcore.DefaultLineEnding,
	}
	options := options{
		console: false,
		level:   "info",
		encoder: defaultEncoderCfg,
	}
	for _, o := range opts {
		o(&options)
	}
	level, err := zapcore.ParseLevel(options.level)
	if err != nil {
		return nil, err
	}
	writerSyncs := make([]zapcore.WriteSyncer, 0, 2)
	fileWriter, err := createFileWriter(options.dir, options.fileName, options.maxAge)
	if err != nil {
		return nil, err
	}
	writerSyncs = append(writerSyncs, zapcore.AddSync(fileWriter))
	if options.console {
		writerSyncs = append(writerSyncs, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(options.encoder),
		zapcore.NewMultiWriteSyncer(writerSyncs...),
		zap.NewAtomicLevelAt(level))
	zapOpts := []zap.Option{
		// zap.AddStacktrace(zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(3),
	}
	if len(options.zapOpts) > 0 {
		zapOpts = append(zapOpts, options.zapOpts...)
	}
	zaplogger := zap.New(core, zapOpts...)
	wraper := &zapLogWraper{
		log: zaplogger,
	}
	return wraper, nil
}

func (w *zapLogWraper) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		w.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
	}
	switch level {
	case log.LevelDebug:
		w.log.Debug("", data...)
	case log.LevelInfo:
		w.log.Info("", data...)
	case log.LevelWarn:
		w.log.Warn("", data...)
	case log.LevelError:
		w.log.Error("", data...)
	}
	return nil
}
