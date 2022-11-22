package weixin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	Http "net/http"
	"strings"
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func GetUserPhone(ctx context.Context, appID, code string) (PhoneWrapper, error) {
	accessToken, err := GetToken(ctx, appID)
	if err != nil {
		log.Infof("get token err: %s.", err)
		return PhoneWrapper{}, err
	}
	url := cstUserPhone + accessToken
	data := map[string]interface{}{
		"code": code,
	}
	bytesData, err := json.Marshal(data)
	if err != nil {
		return PhoneWrapper{}, err
	}
	req, _ := Http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")

	client := &Http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Infof("get user phone from wx service err: %s", err)
		return PhoneWrapper{}, err
	}
	replyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Infof("read resp body err: %s", err)
		return PhoneWrapper{}, err
	}
	reply := PhoneWrapper{}
	if err = json.Unmarshal(replyByte, &reply); err != nil {
		log.Infof("unmarshal resp err: %s.", err)
		return reply, err
	}
	return reply, nil
}

func GetSessionByCode(ctx context.Context, appID string, code string) (Session, error) {
	nonce := util.RandString(20)
	timestamp := time.Now().Unix()
	req := map[string]interface{}{
		"appid":     appID,
		"jsCode":    code,
		"nonce":     nonce,
		"timestamp": timestamp,
	}
	req["sign"] = genYHQOpenSign(req, "kM2xtPkDFeT6qxyO")
	params := make([]string, 0, len(req))
	for k, v := range req {
		params = append(params, fmt.Sprintf("%s=%v", k, v))
	}
	rsp := GetSessionRsp{}
	client, err := http.NewClient(context.Background(), http.WithEndpoint(cstSaleWXProxyDomain))
	if err != nil {
		return Session{}, err
	}
	log.Infof("req:%s", strings.Join(params, "&"))
	err = client.Invoke(context.Background(), "GET",
		"/wechat/open/applet/getLoginSessionInfo?"+strings.Join(params, "&"), nil, &rsp)
	if err != nil {
		log.Errorf("invoke err:%s", err)
	}
	if rsp.Code != 0 {
		return Session{}, fmt.Errorf("code:%d msg:%s", rsp.Code, rsp.Msg)
	}
	return rsp.Data, nil
}

func GetNativeSessionByCode(ctx context.Context, appID, appSecret string, code string) (Session, error) {
	req := map[string]interface{}{
		"appid":      appID,
		"secret":     appSecret,
		"grant_type": "authorization_code",
		"js_code":    code,
	}
	params := make([]string, 0, len(req))
	for k, v := range req {
		params = append(params, fmt.Sprintf("%s=%v", k, v))
	}
	out := Session{}
	rsp := WXSession{}
	client, err := http.NewClient(context.Background(), http.WithEndpoint(cstWeixinApiDomain))
	if err != nil {
		return out, err
	}
	log.Infof("req:%s", strings.Join(params, "&"))
	err = client.Invoke(context.Background(), "GET",
		"/sns/jscode2session?"+strings.Join(params, "&"), nil, &rsp)
	if err != nil {
		log.Errorf("invoke err:%s", err)
	}
	if rsp.Errcode != 0 {
		return out, fmt.Errorf("code:%d msg:%s", rsp.Errcode, rsp.Errmsg)
	}
	out.Openid = rsp.Openid
	out.Unionid = rsp.Unionid
	out.SessionKey = rsp.SessionKey
	return out, nil
}
