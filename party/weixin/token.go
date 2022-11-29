package weixin

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
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
	err := get(addr, data, reply)
	if err != nil {
		return nil, err
	}
	log.Info("token refresh:", addr, reply, err)
	if reply.Code != 0 {
		return nil, fmt.Errorf("token get fail")
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
