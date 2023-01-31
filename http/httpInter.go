package http

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

type HttpInter interface {
	SetMethod(string) HttpInter
	SetURL(string) HttpInter
	SetQueryArgs(interface{}) HttpInter
	SetBody(interface{}) HttpInter
	SetHeader(map[string]string) HttpInter
	SetTimeOut(int64) HttpInter
	Send(context.Context, *zap.SugaredLogger) ([]byte, error)
}

var (
	_ HttpInter = (*baseHttp)(nil)
)

func (b *baseHttp) SetURL(url string) HttpInter {
	b.url = url
	return b
}

func (b *baseHttp) SetMethod(method string) HttpInter {
	b.method = method
	return b
}

func (b *baseHttp) SetQueryArgs(args interface{}) HttpInter {
	b.queryArgs = args
	return b
}

func (b *baseHttp) SetBody(requestData interface{}) HttpInter {
	b.requestData = requestData
	return b
}

func (b *baseHttp) SetHeader(headers map[string]string) HttpInter {
	b.headers = headers
	return b
}
func (b *baseHttp) SetTimeOut(timeOut int64) HttpInter {
	b.timeOut = timeOut
	return b
}

func (b *baseHttp) Send(ctx context.Context, logger *zap.SugaredLogger) ([]byte, error) {
	defer b.reset()
	if b.method == "" {
		b.method = http.MethodPost
	}
	if len(b.headers) == 0 {
		b.headers = make(map[string]string)
	}
	return b.HandleFunc(ctx, logger, b.timeOut, b.method, b.url, b.queryArgs, b.requestData, b.headers)
}

func (b *baseHttp) reset() {
	b.url = ""
	b.method = ""
	b.queryArgs = nil
	b.requestData = nil
	b.headers = make(map[string]string)
}
