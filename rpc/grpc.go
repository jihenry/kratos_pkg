package rpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var (
	pvDiscovery registry.Discovery
	pv          sync.Mutex
	conns       sync.Map
)

func SetDiscovery(r registry.Discovery) {
	pv.Lock()
	pvDiscovery = r
	pv.Unlock()
}

type Option func(*options)

type options struct {
	discovery registry.Discovery
	endpoint  string
	timeout   time.Duration //连接超时时间
}

func WithDiscovery(r registry.Discovery) Option {
	return func(o *options) {
		o.discovery = r
	}
}

func WithEndpoint(endpoint string) Option {
	return func(o *options) {
		o.endpoint = endpoint
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

func Conn(name string, opts ...Option) (*grpc.ClientConn, error) {
	conn, ok := conns.Load(name)
	if ok {
		if client, _ := conn.(*grpc.ClientConn); client != nil && client.GetState() == connectivity.Ready {
			return client, nil
		} else if client != nil {
			client.Close()
		}
		conns.Delete(name)
	}
	nconn, err := doConnect(name, opts...)
	if err != nil {
		return nil, err
	}
	conns.Store(name, nconn)
	return nconn, nil
}

func doConnect(name string, opts ...Option) (*grpc.ClientConn, error) {
	options := options{
		discovery: pvDiscovery,
		timeout:   0,
	}
	for _, o := range opts {
		o(&options)
	}
	dialOpts := []transgrpc.ClientOption{
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
		transgrpc.WithDiscovery(options.discovery),
		transgrpc.WithTimeout(options.timeout),
	}
	var endpoint = options.endpoint
	if options.endpoint != "" {
		dialOpts = append(dialOpts, transgrpc.WithEndpoint(options.endpoint))
	} else {
		endpoint := fmt.Sprintf("discovery:///%s.grpc", name)
		dialOpts = append(dialOpts, transgrpc.WithEndpoint(endpoint))
		dialOpts = append(dialOpts, transgrpc.WithDiscovery(options.discovery))
	}
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		dialOpts...,
	)
	if err != nil {
		log.Error(endpoint, "connect failed:", err.Error())
		return nil, err
	}
	return conn, nil
}
