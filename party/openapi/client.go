package openapi

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type OpenApi interface {
	NBCBSend(ctx context.Context, reqPath string, reqParam interface{}, rspData interface{}) error
}

type OpenApiOption func(*options)

type options struct {
	serverUrl   string //服务器host，优先使用host
	serviceName string //如果使用nacos，则要传递服务名
	discovery   registry.Discovery
}

func WithServerUrl(serverUrl string) OpenApiOption {
	return func(opts *options) {
		opts.serverUrl = serverUrl
	}
}

func WithServerName(serviceName string) OpenApiOption {
	return func(opts *options) {
		opts.serviceName = serviceName
	}
}

func WithDiscovery(discovery registry.Discovery) OpenApiOption {
	return func(opts *options) {
		opts.discovery = discovery
	}
}

type openApiImpl struct {
	client *http.Client
}

var _ OpenApi = (*openApiImpl)(nil)

func NewOpenApiClient(opts ...OpenApiOption) (OpenApi, error) {
	options := options{}
	for _, opt := range opts {
		opt(&options)
	}
	httpOpts := []http.ClientOption{}
	if options.serverUrl != "" {
		httpOpts = append(httpOpts, http.WithEndpoint(options.serverUrl))
	} else if options.serviceName != "" {
		endpoint := fmt.Sprintf("discovery:///%s.http", options.serviceName)
		httpOpts = append(httpOpts, http.WithEndpoint(endpoint))
		httpOpts = append(httpOpts, http.WithDiscovery(options.discovery))
	}
	client, err := http.NewClient(context.Background(), httpOpts...)
	if err != nil {
		return nil, err
	}
	return &openApiImpl{client: client}, nil
}
