package weixin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"gitlab.yeahka.com/gaas/pkg/util"
)

func tokenSaleWellRefresh(ctx context.Context, in *RefreshTokenRequest) (*SalewellToken, error) {
	nonce := util.RandString(20)
	timestamp := time.Now().Unix()
	data := map[string]interface{}{
		"appid":     in.Header.AppId,
		"nonce":     nonce,
		"timestamp": timestamp,
	}
	data["sign"] = genYHQOpenSign(data, "kM2xtPkDFeT6qxyO")
	reply := &SalewellTokenReply{}
	addr := "https://oauth.salewell.com/wechat/open/app/getToken"
	client, err := http.NewClient(context.Background(), http.WithEndpoint("https://oauth.salewell.com"))
	if err != nil {
		return reply.Data, err
	}

	params := make([]string, 0, len(data))
	for k, v := range data {
		params = append(params, fmt.Sprintf("%s=%v", k, v))
	}

	err = client.Invoke(context.Background(), "GET", "/wechat/open/app/getToken?"+strings.Join(params, "&"), nil, &reply)
	if err != nil {
		return reply.Data, err
	}

	log.Info("token refresh:", addr, reply, err)
	if reply.Code != 0 {
		return nil, fmt.Errorf("token get fail")
	}
	return reply.Data, nil
}

func GetToken(ctx context.Context, appID string) (string, error) {
	req := &RefreshTokenRequest{
		Header:  &APIHeader{AppId: appID},
		Request: &RefreshTokenRequest_Request{Force: false},
	}
	reply, err := tokenSaleWellRefresh(ctx, req)
	if err != nil {
		return "", err
	}
	return reply.AccessToken, nil
}
