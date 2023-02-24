package log

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type loggerContextKey struct {
}

func NewContext(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

func FromContext(ctx context.Context, kvs ...interface{}) *log.Helper {
	ctxValueLogger := FromLogger(ctx, kvs...)
	return log.NewHelper(ctxValueLogger)
}

func FromLogger(ctx context.Context, kvs ...interface{}) log.Logger {
	logger := log.GetLogger()
	if ctxLogger, ok := ctx.Value(loggerContextKey{}).(log.Logger); ok {
		logger = ctxLogger
	}
	ctxValueLogger := log.WithContext(ctx, logger)
	if len(kvs) > 0 {
		ctxValueLogger = log.With(ctxValueLogger, kvs...)
	}
	return ctxValueLogger
}
