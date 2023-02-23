package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
)

type RequestData struct {
	RequestId string
}

type requestDataContextKey struct{}

func NewRequestDataContext(ctx context.Context, data *RequestData) context.Context {
	return context.WithValue(ctx, requestDataContextKey{}, data)
}

func FromRequestDataContext(ctx context.Context) *RequestData {
	if ctxData, ok := ctx.Value(requestDataContextKey{}).(*RequestData); ok {
		return ctxData
	}
	return &RequestData{}
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

func GetRequestID(ctx context.Context) string {
	rd := FromRequestDataContext(ctx)
	if rd != nil {
		return rd.RequestId
	}
	return ""
}
