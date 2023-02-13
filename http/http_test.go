package http

import (
	"context"
	"testing"

	klog "github.com/go-kratos/kratos/v2/log"
	"gitlab.yeahka.com/gaas/pkg/log"

	jsoniter "github.com/json-iterator/go"
)

func TestMain(m *testing.M) {
	InitHttpClient(
		WithTimeout(int64(5)),
		WithIdleConnTimeout(60),
		WithMaxIdleConnsPerHost(20),
		WithMaxIdleConns(20),
	)
	log.InitZapLogger(log.ZapLoggerConf{
		Level:       "info",
		FileName:    "test",
		FilePath:    "./http",
		MaxSize:     100,
		MaxBackups:  30,
		MaxAge:      30,
		Compress:    true,
		ShowConsole: true,
	})
	m.Run()
}

func TestBaseHttp_Send(t *testing.T) {
	data, _ := jsoniter.Marshal(map[string]interface{}{
		"channel":  1,
		"gameId":   28,
		"serverId": "cag6hhjurlckotcv54g0",
		"code":     "EC0A2E01-5B9C-43E3-83E7-164ECB5244E212",
		"extend": map[string]interface{}{
			"channelExt": map[string]string{
				"shareFrom":    "",
				"shareFromAct": "",
			},
		},
	})
	ctx := context.Background()
	data, err := NewHttp().SetURL("https://d2-gmapi.yeahkagame.com/gaas/user/login").
		SetMethod("POST").SetBody(data).Send(ctx, klog.GetLogger())
	t.Log("data", string(data), "err", err)
}
