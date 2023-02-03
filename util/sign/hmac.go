package sign

import (
	"fmt"
	"sort"
	"strings"

	ji "github.com/json-iterator/go"
)

//签名秘钥分为未登录秘钥和已登录秘钥。
//未登录秘钥：通过请求参数计算获得。公式为：signKey=upperCase(hex(md5(timestamp + "#" + nonce)))
//已登录秘钥：通过请求参数和登录返回的签名盐计算得到，若返回签名盐为空时，计算方式与未登录秘钥计算方式一致。
//签名盐公式为：signKey=upperCase(hex(md5(timestamp + "#" + nonce + "#" + saltKey)))
/**
1.按ASCII码从小到大排序
2.统一使用UTF8进行编码签名，防止编码方式或特殊字符不兼容问题
3.签名原始串中，字段名和字段值都采用原始值，不进行URL Encode
4.注意整形、浮点型数据参与签名方式（如：浮点数3.10体现为3.1、0.0体现为0）
5.内嵌JSON或ARRAY解析拼接需转字符串且按紧凑方式，即内嵌各K/V或值之间不应有空格或换行符等等
*/
type (
	Sign interface {
		// CreateSaltKeyByUserAgent 创建盐,依赖用户的请求的 user agent
		/**
		@Param userAgent 用户请求的浏览器信息
		@Param sessionId 请求会话
		*/
		CreateSaltKeyByUserAgent(sessionId, userAgent string) Sign
		// BuildSignQueryStr 进行参数签名 返回 返回签名的结果和重放的字典结构
		// ASCII码从小到大排序，空键/值不参与组串
		/**
		@Param data json 序列化的参数
		@return ss string  服务器签名  replayMaps 服务端防重复的参数
		*/
		BuildSignQueryStr(data []byte) (ss string, replayMaps map[string]interface{})
	}
	sign struct {
		saltKey   string
		hashKey   string
		signSlice []string
	}
)

var (
	_            Sign = (*sign)(nil)
	jsonIterator      = ji.Config{
		EscapeHTML:                    false,
		ObjectFieldMustBeSimpleString: true,
		UseNumber:                     true,
	}.Froze()
)

func InitSign(signSlice ...string) Sign {
	return &sign{
		signSlice: signSlice,
	}
}

func (s *sign) splitUserAgent(userAgent string) string {
	agent := strings.Split(userAgent, " ")
	agents := make([]string, 0, 0)
	for _, v := range agent {
		if strings.Contains(v, "Process") {
			continue
		}
		agents = append(agents, v)
	}
	return strings.Join(agents, " ")
}
func (s *sign) CreateSaltKeyByUserAgent(sessionId, userAgent string) Sign {
	s.saltKey = Md5(sessionId + s.splitUserAgent(userAgent))
	return s
}

func (s *sign) BuildSignQueryStr(data []byte) (ss string, replayMaps map[string]interface{}) {
	var (
		maps   = make(map[string]interface{})
		slices = make([]interface{}, 0, len(s.signSlice))
		hk     string
	)
	if err := jsonIterator.Unmarshal(data, &maps); err != nil {
		return
	}
	keys := make([]string, 0, len(maps))
	for _, v := range s.signSlice {
		if value, ok := maps[v]; ok {
			slices = append(slices, value)
		}
	}
	if len(slices) > 0 {
		var (
			buf strings.Builder
		)
		for _, v := range slices {
			val := ToString(v)
			if buf.Len() > 0 {
				buf.WriteByte('#')
			}
			buf.WriteString(val)
		}
		buf.WriteByte('#')
		buf.WriteString(s.saltKey)
		hk = strings.ToUpper(Md5(buf.String()))
	}
	for k := range maps {
		keys = append(keys, k)
	}
	// ASCII码从小到大排序，空键/值不参与组串
	sort.Strings(keys)
	var (
		buf strings.Builder
	)
	for _, key := range keys {
		if value, ok := maps[key]; ok {
			val := ToString(value)
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(key)
			buf.WriteByte('=')
			buf.WriteString(val)
		}
	}
	if len(hk) > 0 {
		buf.WriteByte('&')
		buf.WriteString("key")
		buf.WriteByte('=')
		buf.WriteString(hk)
	}
	ss = Hmac(buf.String(), hk)
	replayMaps = maps
	return
}

func ToString(value interface{}) (vv string) {
	switch value.(type) {
	case string:
		vv = fmt.Sprintf("%s", value)
	default:
		vv, _ = jsonIterator.MarshalToString(value)
	}
	return
}
