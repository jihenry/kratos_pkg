package db

import (
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Option func(*options)

type options struct {
	source      string
	maxConn     int
	maxIdleConn int
	maxLifeTime time.Duration
	logger      logger.Interface
	logLevel    logger.LogLevel
}

func WithSource(source string) Option {
	return func(o *options) {
		o.source = source
	}
}

func WithLogger(logger logger.Interface) Option {
	return func(o *options) {
		o.logger = logger
	}
}

func WithMaxConn(maxConn int) Option {
	return func(o *options) {
		o.maxConn = maxConn
	}
}

func WithMaxIdleConn(maxIdleConn int) Option {
	return func(o *options) {
		o.maxIdleConn = maxIdleConn
	}
}

func WithMaxLifeTime(maxLifeTime time.Duration) Option {
	return func(o *options) {
		o.maxLifeTime = maxLifeTime
	}
}

func WithLogLevel(logLevel logger.LogLevel) Option {
	return func(o *options) {
		o.logLevel = logLevel
	}
}

func NewMysqlClient(opts ...Option) (*gorm.DB, error) {
	options := options{
		maxConn:     100,
		maxIdleConn: 10,
		maxLifeTime: time.Duration(300) * time.Second,
		logger:      nil,
		logLevel:    logger.Silent,
	}
	for _, o := range opts {
		o(&options)
	}
	if util.IsNil(options.logger) {
		options.logger = logger.Default.LogMode(options.logLevel)
	}
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: options.source}),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: options.logger,
		})
	if err != nil {
		return nil, err
	}
	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sdb.SetMaxOpenConns(options.maxConn)
	sdb.SetMaxIdleConns(options.maxIdleConn)
	sdb.SetConnMaxLifetime(options.maxLifeTime)
	return db, nil
}
