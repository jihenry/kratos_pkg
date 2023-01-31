package zaplog

import (
	"go.uber.org/zap"
	glg "gorm.io/gorm/logger"
)

var (
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
)

func InitGormLogger(level int) glg.Interface {
	var (
		l glg.Interface
	)
	if logger != nil {
		l = NewGormLogger(logger)
		l = l.(*GormLogger).LogMode(glg.LogLevel(level))
	} else {
		//走默认的mode
		l = glg.Default.LogMode(glg.LogLevel(level))
	}
	return l
}

// InitZapLogger 初始化zap 日志
func InitZapLogger(zc ZapLoggerConf) {
	sugarLogger = newZapLogger(zc)
	logger = sugarLogger.Desugar()
}

//LoggerWith 重新创建新的日志
func LoggerWith(logger *zap.SugaredLogger, withs ...interface{}) *zap.SugaredLogger {
	return logger.With(withs...)
}
