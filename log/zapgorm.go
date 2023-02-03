package log

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	glg "gorm.io/gorm/logger"
)

type GormLogger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  glg.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

type (
	gormInter glg.Interface
)

var (
	_ gormInter = (*GormLogger)(nil)
)

// NewGormLogger 初始化日志
func NewGormLogger(zl *zap.Logger) *GormLogger {
	return &GormLogger{
		ZapLogger:                 zl,
		LogLevel:                  glg.Warn,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          true,
		IgnoreRecordNotFoundError: false,
	}
}

func (gl *GormLogger) LogMode(l glg.LogLevel) glg.Interface {
	gl.LogLevel = l
	return gl
}

func (gl *GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if gl.LogLevel < glg.Info {
		return
	}
	//输出日志
	gl.ZapLogger.Sugar().Debugf(str, args...)
}

func (gl *GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if gl.LogLevel < glg.Warn {
		return
	}
	//输出日志
	gl.ZapLogger.Sugar().Warnf(str, args...)
}

func (gl *GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if gl.LogLevel < glg.Error {
		return
	}
	//输出日志
	gl.ZapLogger.Sugar().Warnf(str, args...)
}

func (gl *GormLogger) Trace(ctx context.Context, begin time.Time, fun func() (string, int64), err error) {
	if gl.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fun()
	var (
		rowInterface interface{} = rows
	)
	if rows == -1 {
		rowInterface = "-"
	}
	ms := float64(elapsed.Nanoseconds()) / 1e6
	switch {
	// error
	case err != nil && gl.LogLevel >= glg.Error && (!gl.IgnoreRecordNotFoundError || !errors.Is(err, glg.ErrRecordNotFound)):
		gl.ZapLogger.Sugar().Errorf("\n[%.3fms] [rows:%v] %s",
			ms,
			rowInterface,
			sql,
		)
	case gl.LogLevel == glg.Info:
		gl.ZapLogger.Sugar().Infof("\n[%.3fms] [rows:%v] %s",
			ms,
			rowInterface,
			sql,
		)
	case gl.SlowThreshold != 0 && elapsed > gl.SlowThreshold && gl.LogLevel >= glg.Warn:
		gl.ZapLogger.Sugar().Warnf("%s\n[%.3fms] [rows:%v] %s",
			zap.Duration("SLOW SQL >= %v", gl.SlowThreshold),
			ms,
			rowInterface,
			sql,
		)
	}
}
