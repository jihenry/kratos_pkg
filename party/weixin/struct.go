package weixin

const (
	cstSaleWXProxyDomain = "https://oauth.salewell.com" //消费云微信代理平台
	cstWeixinApiDomain   = "https://api.weixin.qq.com"  //微信平台api
	cstUserPhone         = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token="
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

type RefreshTokenRequest struct {
	Header  *APIHeader                   `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Request *RefreshTokenRequest_Request `protobuf:"bytes,2,opt,name=request,proto3" json:"request,omitempty"`
}
type SalewellJsCodeReply struct {
	Code int32                `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string               `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data *SalewellSessionInfo `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}
type RefreshTokenRequest_Request struct {
	Force bool `protobuf:"varint,1,opt,name=force,proto3" json:"force,omitempty"`
}
type SalewellSessionInfo struct {
	Openid     string `protobuf:"bytes,1,opt,name=openid,proto3" json:"openid,omitempty"`
	SessionKey string `protobuf:"bytes,2,opt,name=sessionKey,proto3" json:"sessionKey,omitempty"`
	Unionid    string `protobuf:"bytes,3,opt,name=unionid,proto3" json:"unionid,omitempty"`
}
type APIHeader struct {
	AppId string `protobuf:"bytes,1,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
}
type SalewellTokenReply struct {
	Code int32          `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string         `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data *SalewellToken `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

type SalewellToken struct {
	Appid       string `protobuf:"bytes,1,opt,name=appid,proto3" json:"appid,omitempty"`
	AccessToken string `protobuf:"bytes,2,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
	ExpireTime  string `protobuf:"bytes,3,opt,name=expireTime,proto3" json:"expireTime,omitempty"`
}

type PhoneWrapper struct {
	PhoneInfo PhoneInfo `json:"phone_info"`
}

type PhoneInfo struct {
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}
