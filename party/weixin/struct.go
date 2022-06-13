package weixin

const (
	cstSaleWXProxyDomain = "https://oauth.salewell.com" //消费云微信代理平台
	cstWeixinApiDomain   = "https://api.weixin.qq.com"  //微信平台api
)

type GetSessionRsp struct {
	Code int32   `json:"code,omitempty"`
	Msg  string  `json:"msg,omitempty"`
	Data Session `json:"data,omitempty"`
}

type Session struct {
	Openid     string `json:"openid,omitempty"`
	SessionKey string `json:"sessionKey,omitempty"`
	Unionid    string `json:"unionid,omitempty"`
}

type WXSession struct {
	WXErr
	Openid     string `json:"openid,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
	Unionid    string `json:"unionid,omitempty"`
}

type WXErr struct {
	Errcode int32
	Errmsg  string
}
