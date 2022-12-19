package openapi

import (
	"context"
	"testing"
)

func TestNBCBSend(t *testing.T) {
	client, err := NewOpenApiClient(WithServerUrl("https://t-wxservice.yeahkagame.com"))
	if err != nil {
		t.Errorf("NewOpenApiClient err:%s", err)
		return
	}
	reqParam := struct {
		UnionId string `json:"unionId"`
		Phone   int64  `json:"phone"`
	}{
		UnionId: "oAW4Rs7BSrnHzQtXMr5oob0U13_k",
		Phone:   18844475022,
	}
	rspData := struct {
		Promoter   bool `json:"promoter"`
		CreditCard bool `json:"creditCard"`
		Grouping   bool `json:"grouping"`
	}{
		Promoter: true,
	}
	err = client.NBCBSend(context.Background(), "/town/queryTag", reqParam, &rspData)
	if err != nil {
		t.Errorf("NBCBSend err:%s", err)
	}
}
