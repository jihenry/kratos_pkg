package middleware

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"

	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/go-kratos/kratos/v2/errors"

	"github.com/go-kratos/kratos/v2/transport"

	"github.com/go-kratos/kratos/v2/metadata"

	uuid "github.com/satori/go.uuid"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"gitlab.yeahka.com/gaas/pkg/zaplog"
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
					zaplog.FromContext(ctx).Errorf("%v: %+v\n%s\n", rerr, req, buf)
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
			zaplog.FromContext(ctx).Infof("%s kind:%s | component:%s | operation: %s | args:%s | code:%d | reason:%s | stack:%s | latency:%f ",
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

const (
	RequestID = `X-Request-Id`
	XTraceID  = `X-Trace-Id`
	XSpanId   = `X-Span-Id`
)

func Value(ctx context.Context, v interface{}) interface{} {
	if v, ok := v.(log.Valuer); ok {
		return v(ctx)
	}
	return v
}

func (m *middle) RequestContext() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			//从上游获取到请求的 X-Request-Id]
			var (
				requestId string
			)
			if val, ok := GetMetaData(ctx, "X-Request-Id")["X-Request-Id"]; ok {
				requestId = val
			}
			if requestId == "" {
				requestId = UUId()
			}
			newCtx := util.NewRequestContext(ctx, util.RequestData{"requestId": requestId})
			//日志的上下文
			logCtx := zaplog.NewLoggerContext(newCtx,
				zaplog.LoggerWith(zaplog.FromContext(newCtx), []interface{}{
					RequestID, requestId,
					XTraceID, Value(newCtx, tracing.TraceID()),
					XSpanId, Value(newCtx, tracing.SpanID()),
				}...,
				),
			)
			return handler(logCtx, req)
		}
	}
}

// SetMetaData 设置请求的RPC请求的上下文传递参数
func SetMetaData(ctx context.Context, metaData map[string]string) context.Context {
	var (
		globalKey = `x-md-global-`
	)
	for k := range metaData {
		var (
			extra string
		)
		if md, ok := metadata.FromServerContext(ctx); ok {
			extra = md.Get(globalKey + k)
		}
		if extra == "" {
			if val, ok := metaData[k]; ok {
				extra = val
			}
		}
		ctx = metadata.AppendToClientContext(ctx, globalKey+k, extra)
	}
	return ctx
}

func GetMetaData(ctx context.Context, key ...string) map[string]string {
	var (
		globalKey = `x-md-global-`
		ans       = make(map[string]string)
	)
	var (
		extra string
	)
	for k := range key {
		val := key[k]
		if tr, ok := transport.FromServerContext(ctx); ok {
			extra = tr.RequestHeader().Get(globalKey + strings.ToLower(val))
		}
		if tr, ok := transport.FromClientContext(ctx); ok {
			extra = tr.ReplyHeader().Get(globalKey + strings.ToLower(val))
		}
		if md, ok := metadata.FromServerContext(ctx); ok {
			extra = md.Get(globalKey + val)
		}
		if md, ok := metadata.FromClientContext(ctx); ok {
			extra = md.Get(globalKey + val)
		}
		ans[val] = extra
	}
	return ans
}

func (m *middle) extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func (m *middle) extractError(err error) (zaplog.Level, string) {
	if err != nil {
		return zaplog.LevelError, fmt.Sprintf("%+v", err)
	}
	return zaplog.LevelInfo, ""
}

func UUId() string {
	return uuid.NewV4().String()
}
