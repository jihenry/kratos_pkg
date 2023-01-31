package nacos

import (
	"testing"

	"gitlab.yeahka.com/gaas/pkg/zaplog"
)

func TestInitServer(t *testing.T) {
	InitNacosClient(NacosConfig{
		NamespaceID:         "public",
		RotateTime:          "1h",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
		CacheDir:            "../nacos/cache",
		LogDir:              "../nacos/log",
		MaxAge:              3,
		NacosServer: ServerConfig{
			IP:   "127.0.0.1",
			Port: 8848,
		},
	}, []string{"account"}...)
}

func TestGetConn(t *testing.T) {
	zaplog.InitZapLogger(zaplog.ZapLoggerConf{
		Level:       "info",
		FileName:    "server",
		FilePath:    "../zaplog/",
		MaxSize:     100,
		MaxBackups:  30,
		MaxAge:      30,
		Compress:    true,
		ShowConsole: true,
	})
	InitNacosClient(NacosConfig{
		NamespaceID:         "public",
		RotateTime:          "1h",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "debug",
		CacheDir:            "../nacos/cache",
		LogDir:              "../nacos/log",
		MaxAge:              3,
		NacosServer: ServerConfig{
			IP:   "127.0.0.1",
			Port: 8848,
		},
	}, []string{"account"}...)
}
