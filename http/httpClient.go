package http

import (
	"net/http"
	"sync"
	"time"
)

var (
	once       sync.Once
	httpClient *http.Client
)

type Option func(*option)

type option struct {
	timeout             int64
	maxIdleConns        int
	maxIdleConnsPerHost int
	idleConnTimeout     int
}

func WithTimeout(timeout int64) Option{
	return func(o *option){
		o.timeout = timeout
	}
}

func WithMaxIdleConns(maxIdleConns int) Option{
	return func(o *option){
		o.maxIdleConns = maxIdleConns
	}
}

func WithMaxIdleConnsPerHost(maxIdleConnsPerHost int) Option{
	return func(o *option){
		o.maxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

func WithIdleConnTimeout(idleConnTimeout int) Option{
	return func(o *option){
		o.idleConnTimeout = idleConnTimeout
	}
}

func newHttpClient(options ...Option) *http.Client {
	opt := new(option)
	for _, o := range options{
		o(opt)
	}
	return &http.Client{
		Timeout: time.Duration(opt.timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        opt.maxIdleConns,
			MaxIdleConnsPerHost: opt.maxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(opt.idleConnTimeout) * time.Second,
		},
	}
}

// InitHttpClient  初始化Http服务
func InitHttpClient(options ...Option) *http.Client {
	once.Do(func() {
		httpClient = newHttpClient(options...)
	})
	return httpClient
}
