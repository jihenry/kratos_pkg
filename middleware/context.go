package middleware

import "context"

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
