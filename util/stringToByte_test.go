package util

import (
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

type takeBoxReq struct {
	FriendID uint64 `json:"fid" binding:"required"` //好友id
}

func TestNewPointer(t *testing.T) {
	data := `{"fid":"20382925753155624"}`
	req := new(takeBoxReq)
	jsoniter.Unmarshal([]byte(data), &req)
	fmt.Println(data)
}

func BenchmarkNewPointer(b *testing.B) {
	b.ReportAllocs()
	str1 := `{"idField":"userUuid","idValue":"opXJY6PmordKyu-LoSuFhLHstELQ","once":1658221489,"nickName":"中文测试"}`
	for i := 0; i < b.N; i++ {
		bytes := NewPointer().StringToByte(str1)
		_ = NewPointer().ByteToString(bytes)
	}
}

func BenchmarkStringToByte(b *testing.B) {
	b.ReportAllocs()
	str1 := `{"idField":"userUuid","idValue":"opXJY6PmordKyu-LoSuFhLHstELQ","once":1658221489,"nickName":"中文测试"}`
	for i := 0; i < b.N; i++ {
		bytes := []byte(str1)
		_ = string(bytes)
	}
}
