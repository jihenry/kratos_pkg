package sign

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

func Hmac(data, key string) string {
	hm := hmac.New(md5.New, []byte(key))
	hm.Write([]byte(data))
	return hex.EncodeToString(hm.Sum([]byte("")))
}

func Md5(data string) string {
	md := md5.New()
	md.Write([]byte(data))
	return hex.EncodeToString(md.Sum([]byte("")))
}
