package sign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/url"
	"reflect"
	"sort"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// rsaEncrypt 签名
func rsaEncrypt(plainText []byte, hash crypto.Hash, key string) string {
	//获取私钥
	block, _ := pem.Decode(blockMap[key].PrivateKey)
	publicKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	h := crypto.Hash.New(hash) //进行加密
	h.Write(plainText)
	signature, err := rsa.SignPKCS1v15(rand.Reader, publicKeyInterface.(*rsa.PrivateKey), hash, h.Sum(nil))
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(signature)
}

// rsaDecrypt 签名验证
func rsaDecrypt(data []byte, sign string, hash crypto.Hash, key string) bool {
	block, _ := pem.Decode(blockMap[key].PublicKey)
	public, _ := x509.ParsePKIXPublicKey(block.Bytes)
	bytes, _ := base64.StdEncoding.DecodeString(sign)
	h := crypto.Hash.New(hash) //进行加密
	h.Write(data)
	hashed := h.Sum(nil)
	if err := rsa.VerifyPKCS1v15(public.(*rsa.PublicKey), hash, hashed, bytes); err != nil {
		return false
	}
	return true
}

func structToMap(inter interface{}) map[string]interface{} {
	t := reflect.TypeOf(inter)
	values := make(map[string]interface{}, 0)
	v := reflect.ValueOf(inter)
	for i := 0; i < t.NumField(); i++ {
		keys := t.Field(i)
		str := keys.Name[0]
		key0 := strings.ToLower(string(str))
		key := key0 + keys.Name[1:]
		value := v.Field(i).Interface()
		values[key] = value
	}
	return values
}

// httpBuildQuery ascii码从小到大排序后使用QueryString的格式拼接
func httpBuildQuery(query interface{}) string {
	var (
		maps = make(map[string]interface{})
	)
	switch query.(type) {
	case map[string]interface{}, map[string]string:
		data, _ := jsoniter.Marshal(query)
		jsonIterator.Unmarshal(data, &maps)
	default:
		maps = structToMap(query)
	}
	keys := make([]string, 0, len(maps))
	for k := range maps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	urls := url.Values{}
	for _, v := range keys {
		urls.Add(v, ToString(maps[v]))
	}
	return urls.Encode()
}

// HttpBuildQueryRsa 生成签名
func HttpBuildQueryRsa(query interface{}, key string) string {
	return rsaEncrypt([]byte(httpBuildQuery(query)), crypto.SHA256, key)
}

// HttpBuildRsaDecrypt 签名验签
func HttpBuildRsaDecrypt(query interface{}, sign string, key string) bool {
	return rsaDecrypt([]byte(httpBuildQuery(query)), sign, crypto.SHA256, key)
}
