package openapi

import (
	"context"
	"fmt"
	"time"

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
	timeout     time.Duration
}

func WithServerUrl(serverUrl string) OpenApiOption {
	return func(opts *options) {
		if serverUrl != "" {
			opts.serverUrl = serverUrl
		}
	}
}

func WithServiceName(serviceName string) OpenApiOption {
	return func(opts *options) {
		if serviceName != "" {
			opts.serviceName = serviceName
		}
	}
}

func WithTimeout(timeout time.Duration) OpenApiOption {
	return func(opts *options) {
		if timeout != 0 {
			opts.timeout = timeout
		}
	}
}

func WithDiscovery(discovery registry.Discovery) OpenApiOption {
	return func(opts *options) {
		if discovery != nil {
			opts.discovery = discovery
		}
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
	httpOpts := []http.ClientOption{
		http.WithTimeout(options.timeout),
		// http.WithBlock(),
	}
	if options.serverUrl != "" {
		httpOpts = append(httpOpts, http.WithEndpoint(options.serverUrl))
	} else if options.serviceName != "" {
		endpoint := fmt.Sprintf("discovery:///%s", options.serviceName)
		httpOpts = append(httpOpts, http.WithEndpoint(endpoint))
		httpOpts = append(httpOpts, http.WithDiscovery(options.discovery))
	} else {
		return nil, fmt.Errorf("conf isn't supported")
	}
	client, err := http.NewClient(context.Background(), httpOpts...)
	if err != nil {
		return nil, err
	}
	return &openApiImpl{client: client}, nil
}
