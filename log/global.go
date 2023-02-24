package log

import (
	"sync"

	klog "github.com/go-kratos/kratos/v2/log"
)

var (
	globalLogger klog.Logger
	globalHelper *klog.Helper
	lock         sync.Mutex
)

func SetLogger(logger klog.Logger) {
	lock.Lock()
	defer lock.Unlock()
	globalLogger = logger
	globalHelper = klog.NewHelper(logger)
}

func GetLoggerHelper() *klog.Helper {
	return globalHelper
}

func GetLogger() klog.Logger {
	return globalLogger
}

// Log Print log by level and keyvals.
func Log(level klog.Level, keyvals ...interface{}) {
	globalHelper.Log(level, keyvals...)
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	globalHelper.Debug(a...)
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	globalHelper.Debugf(format, a...)
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...interface{}) {
	globalHelper.Debugw(keyvals...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	globalHelper.Info(a...)
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	globalHelper.Infof(format, a...)
}

// Infow logs a message at info level.
func Infow(keyvals ...interface{}) {
	globalHelper.Infow(keyvals...)
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	globalHelper.Warn(a...)
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	globalHelper.Warnf(format, a...)
}

// Warnw logs a message at warnf level.
func Warnw(keyvals ...interface{}) {
	globalHelper.Warnw(keyvals...)
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	globalHelper.Error(a...)
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	globalHelper.Errorf(format, a...)
}

// Errorw logs a message at error level.
func Errorw(keyvals ...interface{}) {
	globalHelper.Errorw(keyvals...)
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	globalHelper.Fatal(a...)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	globalHelper.Fatalf(format, a...)
}

// Fatalw logs a message at fatal level.
func Fatalw(keyvals ...interface{}) {
	globalHelper.Fatalw(keyvals...)
}

func With(l klog.Logger, kv ...interface{}) klog.Logger {
	return klog.With(l, kv...)
}
