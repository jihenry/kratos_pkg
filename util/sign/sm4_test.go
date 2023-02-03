package sign

import (
	"encoding/base64"
	"testing"
)

func TestNewCipher(t *testing.T) {
	key := []byte("Q6dWAr5gc4IKZPLS")
	in := []byte(`{"idField":"userUuid","idValue":"opXJY6PmordKyu-LoSuFhLHstELQ","once":1658221489,"nickName":"中文测试"}`)

	plainTextWithPadding := PKCS5Padding(in, BlockSize)
	data, err := ECBEncrypt(key, plainTextWithPadding)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(base64.StdEncoding.EncodeToString(data))
	text, _ := base64.StdEncoding.DecodeString("9GSek4vcinF7q2wP7qMoMszbWKsTSYAj3qXM24vvxvN2KNp6TF1SmClhXoNCMv6tho1/DA8ZNRXW6hyQE9avJRn347C/FOMLtEHrDc/1C5xa2ibfY3LyxNOmsmgGKFapyn2akepGvXHVUZjOuL9Q7Q==")
	dec, err := ECBDecrypt(key, text)
	t.Log(string(dec))
}
