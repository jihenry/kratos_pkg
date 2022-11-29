package openapi

import (
	"context"
	"testing"
)

func TestNBCBSend(t *testing.T) {
	client, err := NewOpenApiClient(WithServerUrl("http://10.30.20.207:8083"))
	if err != nil {
		t.Errorf("NewOpenApiClient err:%s", err)
		return
	}
	reqParam := struct {
		UnionId string `json:"unionId"`
		Phone   int64  `json:"phone"`
	}{
		UnionId: "11",
		Phone:   15818799620,
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
