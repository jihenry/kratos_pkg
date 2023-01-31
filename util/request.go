package util

import (
	"context"
)

type RequestData map[string]string
type requestKey struct{}

func NewRequestContext(ctx context.Context, rd RequestData) context.Context {
	return context.WithValue(ctx, requestKey{}, rd)
}

// FromClientRequestContext 获取每个请求 上下文
func FromClientRequestContext(ctx context.Context) string {
	if md, ok := ctx.Value(requestKey{}).(RequestData); ok {
		return md["requestId"]
	}
	return ""
}
