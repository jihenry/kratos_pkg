package middleware

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"gitlab.yeahka.com/gaas/pkg/util"
)

type LoggerOption func(*loggerOption)
type loggerOption struct {
	reply bool
}

func WithReply(reply bool) LoggerOption {
	return func(opt *loggerOption) {
		opt.reply = reply
	}
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}

// Server is an server logging middleware.
func ServerLogger(logger log.Logger, options ...LoggerOption) middleware.Middleware {
	opt := loggerOption{}
	for _, option := range options {
		option(&opt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)
			kv := []interface{}{
				"kind", "server",
				"component", kind,
				"operation", operation,
				"args", extractArgs(req),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Seconds(),
				"trace_id", tracing.TraceID()(ctx),
				"span_id", tracing.SpanID()(ctx),
			}
			if err == nil && opt.reply {
				mo, _ := util.JSON.MarshalToString(reply)
				kv = append(kv, "reply", mo)
			}
			_ = log.WithContext(ctx, logger).Log(level, kv...)
			return
		}
	}
}
