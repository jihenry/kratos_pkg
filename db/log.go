package db

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ logger.Interface = (*gormLogger)(nil)

type gormLogger struct {
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func NewGormLogger(level logger.LogLevel) logger.Interface {
	l := &gormLogger{
		LogLevel:                  level,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          true,
		IgnoreRecordNotFoundError: false,
	}
	return l
}

// Error implements logger.Interface
func (l *gormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	log.Errorf(str, args...)
}

// Info implements logger.Interface
func (l *gormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	log.Infof(str, args...)
}

// LogMode implements logger.Interface
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &gormLogger{
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

// Trace implements logger.Interface
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		var rowInterface interface{} = rows
		if rows == -1 {
			rowInterface = "-"
		}
		log.Errorf("\n[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rowInterface, sql)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		var rowInterface interface{} = rows
		if rows == -1 {
			rowInterface = "-"
		}
		slowLogField := zap.Duration("SLOW SQL >= %v", l.SlowThreshold)
		log.Warnf("%s\n[%.3fms] [rows:%v] %s", slowLogField, float64(elapsed.Nanoseconds())/1e6, rowInterface, sql)

	case l.LogLevel == logger.Info:
		sql, rows := fc()
		var rowInterface interface{} = rows
		if rows == -1 {
			rowInterface = "-"
		}
		log.Infof("\n[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rowInterface, sql)
	}
}

// Warn implements logger.Interface
func (l *gormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	log.Warnf(str, args...)
}
