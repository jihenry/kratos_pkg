package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Option func(*options)

type options struct {
	addr     string
	password string
	db       int
	poolSize int
}

func WithAddr(addr string) Option {
	return func(o *options) {
		o.addr = addr
	}
}

func WithPassword(password string) Option {
	return func(o *options) {
		o.password = password
	}
}

func WithDb(db int) Option {
	return func(o *options) {
		o.db = db
	}
}

func WithPoolSize(poolSize int) Option {
	return func(o *options) {
		o.poolSize = poolSize
	}
}

func NewRedisClient(opts ...Option) (*redis.Client, error) {
	options := options{
		poolSize: 100,
	}
	for _, o := range opts {
		o(&options)
	}
	cli := redis.NewClient(&redis.Options{
		PoolSize: options.poolSize,
		Addr:     options.addr,
		Password: options.password,
		DB:       options.db,
	})
	_, err := cli.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}
	return cli, nil
}
