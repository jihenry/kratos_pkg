package log

import (
	"sync"

	klog "github.com/go-kratos/kratos/v2/log"
)

var (
	global *klog.Helper
	lock   sync.Mutex
)

func SetLogger(logger *klog.Helper) {
	lock.Lock()
	defer lock.Unlock()
	global = logger
}

func GetLogger() *klog.Helper {
	return global
}

// Log Print log by level and keyvals.
func Log(level klog.Level, keyvals ...interface{}) {
	global.Log(level, keyvals...)
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	global.Debug(a...)
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	global.Debugf(format, a...)
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...interface{}) {
	global.Debugw(keyvals...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	global.Info(a...)
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	global.Infof(format, a...)
}

// Infow logs a message at info level.
func Infow(keyvals ...interface{}) {
	global.Infow(keyvals...)
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	global.Warn(a...)
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	global.Warnf(format, a...)
}

// Warnw logs a message at warnf level.
func Warnw(keyvals ...interface{}) {
	global.Warnw(keyvals...)
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	global.Error(a...)
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	global.Errorf(format, a...)
}

// Errorw logs a message at error level.
func Errorw(keyvals ...interface{}) {
	global.Errorw(keyvals...)
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	global.Fatal(a...)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	global.Fatalf(format, a...)
}

// Fatalw logs a message at fatal level.
func Fatalw(keyvals ...interface{}) {
	global.Fatalw(keyvals...)
}
