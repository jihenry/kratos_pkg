package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	privKey = `
-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBALSVojdkqMsh53cI
bHzt5GfPLOb77009WucVS8of1nH0CagB0BNecyxohiSpo2WDEwAc24QPeyoHn59b
gVib0Wi7nhZMq7J8zOoWg9HIHg+uR//9LA77wK0evYJ8YfHVB27InFs1IZhiERK3
CMAqDwWqGZyEj2Yht0YG2COZLdLhAgMBAAECgYEAo5pE4oZxXccTmoWpM+2aVmod
tg5dGM8TQfPLPA1oDMkYznsF9eZF1d/EWAbQH7GGTz3VqmkUHlnVxVvzbUGNjx13
BqUUIPdDh8TuER20QEje1sk7IulWP5DgY+nQLV/65vuDDKUHeGolEQLVrIuaMJDV
aJpuncd8/JPP/1F1EbkCQQDaiIxxmiNbyFjAxPopOPNX5NEXZFu9yjUuIU82Q11d
Jm+B+C7F8FdUGml2xGIKY/qDx0PEKT0WK6LO9zBqThZXAkEA04uBO2M2hL6EI6l7
bhBdf7kDLVhfabH2Cs7VOi81jQfZgQv6763ayKkHeDzreKy5cIIP1WZxqVcI/uS4
2UtthwJAHwYzqg0P5//RWcydFy0WnuvFI2UEATWrxxjDfhiiMI88VV8+hKtSOoZl
Yo8OvBrlfb/URwzztyoKuwcswGrFkQJAMt1iT3NFkpl0kFaaFRbeRG2p8+dB2dou
fN7KqljbmXN/uuW0ipjU+FacMy8Ct1tgo0rCn98oCT2iLhe00pquVQJBAMW5UzZC
/WH6Nr/PxxlcWN2QxTbaxJnQp9GVfQqDSscQXCHJqqCWuYXUVETjKFO5cNH/qarL
kxJY6hZUvab5X5U=
-----END PRIVATE KEY-----
`
	pubKey = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0laI3ZKjLIed3CGx87eRnzyzm
++9NPVrnFUvKH9Zx9AmoAdATXnMsaIYkqaNlgxMAHNuED3sqB5+fW4FYm9Fou54W
TKuyfMzqFoPRyB4Prkf//SwO+8CtHr2CfGHx1QduyJxbNSGYYhEStwjAKg8Fqhmc
hI9mIbdGBtgjmS3S4QIDAQAB
-----END PUBLIC KEY-----
`
)

var (
	encrypt string
	decrypt string
	RSA     = &SHA1withRSA{}
)

func TestMain(m *testing.M) {
	RSA.SetPriKey([]byte(privKey))
	RSA.SetPublicKey([]byte(pubKey))
	m.Run()
}

func TestSign(t *testing.T) {
	sign, _ := RSA.Sign([]byte("nihao"))
	if !RSA.CheckSign(sign, []byte("nihao")) {
		t.Fatal("check sign verify fail")
	}
	now := time.Now()
	fmt.Println(now.UnixNano() / 1e6)
	data := make(map[string]int32)
	data["pageSize"] = 20
	data["currentPage"] = 1
	bb, _ := json.Marshal(data)
	var Req = struct {
		ChannelKey string `json:"channelKey"`
		Data       string `json:"data"`
		// Sign       string `json:"sign"`
		TimeStamp int64 `json:"timeStamp"`
	}{
		ChannelKey: "ck0761892307",
		Data:       string(bb),
		TimeStamp:  time.Now().UnixNano() / 1e6,
	}

	byteData, err := json.Marshal(Req)
	if err != nil {
		log.Errorf("encode json err: %v", err)
		t.Fatalf("Err:%v", err)
	}
	log.Info("sign str: " + string(byteData))
	sig, _ := RSA.Sign(byteData)
	log.Info(sig)
	// Req.Sign = sig
	js, err := json.Marshal(Req)
	if err != nil {
		log.Errorf("encode json err: %v", err)
		t.Fatalf("Err:%v", err)
	}
	out := make(map[string]interface{})
	_ = json.Unmarshal(js, &out)
	out["sign"] = sig
	res, _ := json.Marshal(out)
	fmt.Println("req: " + string(res))
}

func TestEncrypt(t *testing.T) {
	encrypt, _ = RSA.Encrypt([]byte("hello world"))
}

func TestDecrypt(t *testing.T) {
	data, _ := base64.StdEncoding.DecodeString(encrypt)
	decrypt, _ = RSA.Decrypt(data)
}

func TestFlow(t *testing.T) {
	TestEncrypt(t)
	TestDecrypt(t)
}
