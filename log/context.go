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

func FromContext(ctx context.Context, opts ...log.Option) *log.Helper {
	logger := log.GetLogger()
	if ctxLogger, ok := ctx.Value(loggerContextKey{}).(log.Logger); ok {
		logger = ctxLogger
	}
	ctxValueLogger := log.WithContext(ctx, logger)
	return log.NewHelper(ctxValueLogger, opts...)
}
