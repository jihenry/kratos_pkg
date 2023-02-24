package middleware

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"

	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/go-kratos/kratos/v2/errors"

	"github.com/go-kratos/kratos/v2/transport"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"gitlab.yeahka.com/gaas/pkg/log"
)

type (
	Middleware interface {
		Recovery() middleware.Middleware         //异常捕获
		RPCLogging(string) middleware.Middleware //GRPC 日志输出
		RequestContext() middleware.Middleware   //请求的上下文
	}
	middle struct {
		ErrUnknownRequest *errors.Error
		handler           func(ctx context.Context, req, err interface{}) error
	}
)

var (
	_ Middleware = (*middle)(nil)
)

func InitMiddleware() Middleware {
	errs := errors.InternalServer("UNKNOWN", "unknown request error")
	return &middle{
		ErrUnknownRequest: errs,
		handler: func(ctx context.Context, req, err interface{}) error {
			return errs
		},
	}
}

func (m *middle) Recovery() middleware.Middleware {
	m.handler = func(ctx context.Context, req, err interface{}) error {
		return m.ErrUnknownRequest
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if rerr := recover(); rerr != nil {
					buf := make([]byte, 64<<10) //nolint:gomnd
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					log.FromContext(ctx).Errorf("%v: %+v\n%s\n", rerr, req, buf)
					err = m.handler(ctx, req, rerr)
				}
			}()
			return handler(ctx, req)
		}
	}
}

func (m *middle) RPCLogging(template string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				message   string
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
				message = se.Message
			}
			level, stack := m.extractError(err)
			log.FromContext(ctx).Infof("%s kind:%s | component:%s | operation: %s | args:%s | code:%d | reason:%s | stack:%s | latency:%f ",
				level,
				template,
				kind,
				operation,
				m.extractArgs(req),
				code,
				message,
				stack,
				time.Since(startTime).Seconds(),
			)
			return
		}
	}
}

func (m *middle) RequestContext() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			requestId := GetMetaData(ctx, RequestKeyXRequestID)[RequestKeyXRequestID] //从上游获取到请求的 X-Request-Id]
			if requestId == "" {
			}
			requestId = util.New()
			requestDataCtx := NewRequestDataContext(ctx, &RequestData{RequestId: requestId})
			loggerCtx := log.NewContext(requestDataCtx, klog.With(klog.GetLogger(),
				RequestKeyXRequestID, requestId, RequestKeyXTraceID, tracing.TraceID(), RequestKeyXSpanId, tracing.SpanID()))
			return handler(loggerCtx, req)
		}
	}
}

func (m *middle) extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func (m *middle) extractError(err error) (klog.Level, string) {
	if err != nil {
		return klog.LevelError, fmt.Sprintf("%+v", err)
	}
	return klog.LevelInfo, ""
}
