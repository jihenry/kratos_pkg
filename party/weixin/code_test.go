package weixin

import (
	"context"
	"fmt"
	"testing"
)

func TestGetSession(t *testing.T) {
	cases := []struct {
		uri   string
		appid string
		code  string
	}{
		{"水果合合合", "wxe50aaf462a3041b0", "093n3xll25Z0Z84r6vnl2yCX9A2n3xlv"},
		{"神奇板子", "wx3cdaa9114157de80", "093n3xll25Z0Z84r6vnl2yCX9A2n3xlv"},
	}
	for _, ca := range cases {
		rsp, err := GetSessionByCode(context.Background(), ca.appid, ca.code)
		if err != nil {
			t.Errorf("err:%s", err)
		} else {
			fmt.Printf("接口:%s appid:%s rsp:%s\n", ca.uri, ca.appid, rsp)
		}
	}
}

func TestGetNativeSession(t *testing.T) {
	cases := []struct {
		uri       string
		appid     string
		appSecret string
		code      string
	}{
		{"神奇板子", "wx3cdaa9114157de80", "1de4fdbf3ca2634fa4e49adf807e9f51", "093n3xll25Z0Z84r6vnl2yCX9A2n3xlv"},
	}
	for _, ca := range cases {
		rsp, err := GetNativeSessionByCode(context.Background(), ca.appid, ca.appSecret, ca.code)
		if err != nil {
			t.Errorf("err:%s", err)
		} else {
			fmt.Printf("接口:%s appid:%s rsp:%s\n", ca.uri, ca.appid, rsp)
		}
	}
}

func TestGetUserPhone(t *testing.T) {
	phone, _ := GetUserPhone(context.Background(), "wx4e443d9613e66610", "aaaa")
	fmt.Printf("phone: %+v", phone)
}
