package sm4

import (
	"testing"
)

func TestNewECBSm4(t *testing.T) {
	key := []byte("Q6dWAr5gc4lKZPLS")
	//in := []byte(`{"idField":"userUuid","idValue":"opXJY6PmordKyu-LoSuFhLHstELQ","once":1658221489,"nickName":"中文测试"}`)
	//data, err := NewECBSm4(16).PKCS5Padding(in).ECBEncrypt(key)
	//t.Log("data", base64.StdEncoding.EncodeToString(data), "err:", err)
	str := "wsJ7iIT801N3kLDcCgfk0r9Y2ShQo3eiEFzTNF+5aLWtA+S9l+kN6kYrVnkgWdFge+8KCKK29yLADz1rkxvoPNESoYPvRCv0armMGHNeLIyt/MGRomLhsQgnOX1UtxEA3bMsHzQVFg/A1f8FZaAYwy31RvxAb0UVztxHYXP8n5OusgkIoDVJBr0e38XfB3/m"
	dec, _ := NewECBSm4(16).SetCipherText(str).ECBDecrypt(key)
	t.Log(string(dec))
}
