package sign

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// BuildSignQueryStr 签名验证
func TestBuildSignQueryStr(t *testing.T) {
	should := require.New(t)
	ua := `Mozilla/5.0 (Linux; Android 11; RMX3350 Build/RP1A.200720.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/3185 MMWEBSDK/20220105 Mobile Safari/537.36 MMWEBID/9402 MicroMessenger/8.0.19.2080(0x2800133D) Process/appbrand0 WeChat/arm64 Weixin NetType/4G Language/zh_CN ABI/arm64 MiniProgramEnv/android`
	data := `{"ActivityID":"c6o6c5m0desci9hbc21g","timestamp":1646362051391,"nonce":"yeahkagamefi8br8ct57j"}`
	sign, _ := InitSign([]string{"timestamp", "nonce"}...).
		CreateSaltKeyByUserAgent("538fe76d-f299-4e8f-85b5-40665474eab4", ua).
		BuildSignQueryStr([]byte(data))
	should.Equal("1a570d4dfdd29278b4388173c40eeb52f", sign)
}
