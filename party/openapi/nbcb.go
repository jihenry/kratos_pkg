package openapi

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gitlab.yeahka.com/gaas/pkg/util"
)

type nbcbSendHeader struct {
	RqsJrnlNo string `json:"Rqs_Jrnl_No"`
	RspDt     string `json:"Rsp_Dt"`
	RspTm     string `json:"Rsp_Tm"`
}

type NbcbSdkError struct {
	ErrorCode string `json:"errorCode,omitempty"`
	Message   string `json:"message,omitempty"`
	Status    int32  `json:"status,omitempty"`
	ErrorMsg  string `json:"errorMsg,omitempty"`
}

func (e *NbcbSdkError) Error() string {
	return fmt.Sprintf("code:%s, msg:%s", e.ErrorCode, e.ErrorMsg)
}

type nbcbSdkData struct {
	NbcbSdkError
	nbcbSendPathData
}

type nbcbSendPathData struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type nbcbSendRsp struct {
	Head nbcbSendHeader `json:"Head"`
	Data nbcbSdkData    `json:"Data"`
}

func (c *openApiImpl) NBCBSend(ctx context.Context, reqPath string, reqParam interface{}, rspData interface{}) error {
	startTime := time.Now()
	paramJson, err := util.JSON.MarshalToString(reqParam)
	if err != nil {
		return err
	}
	req := map[string]string{
		"serviceID": reqPath,
		"json":      paramJson,
	}
	log.Infof("NBCBSend reqPath:%s paramJson:%s", reqPath, paramJson)
	pathDataType := reflect.TypeOf(rspData)
	if pathDataType.Kind() != reflect.Ptr {
		return fmt.Errorf("data type is not pointer")
	}
	params := make([]string, 0, len(req))
	for k, v := range req {
		params = append(params, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}
	rsp := nbcbSendRsp{
		Data: nbcbSdkData{
			nbcbSendPathData: nbcbSendPathData{
				Data: rspData,
			},
		},
	}
	err = c.client.Invoke(context.Background(), "GET",
		"/nbopen/send?"+strings.Join(params, "&"), nil, &rsp)
	if err != nil {
		log.Errorf("invoke err:%s", err)
		return err
	}
	costTime := time.Since(startTime).Seconds()
	if rsp.Data.NbcbSdkError.ErrorCode != "" {
		log.Errorf("NBCBSend reqPath:%s nbcbSendHeader:%+v NbcbSdkError:%+v costTime:%0.3fs", reqPath, rsp.Head, rsp.Data.NbcbSdkError, costTime)
		return &rsp.Data.NbcbSdkError
	}
	log.Infof("NBCBSend reqPath:%s nbcbSendHeader:%+v nbcbSendPathData:%+v costTime:%0.3fs", reqPath, rsp.Head, rspData, costTime)
	return nil
}
