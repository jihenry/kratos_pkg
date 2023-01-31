package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// HandleFunc
type HandleFunc func(context.Context, *zap.SugaredLogger, int64, string, string, interface{}, interface{}, map[string]string) ([]byte, error)

type baseHttp struct {
	url         string
	method      string
	timeOut     int64
	queryArgs   interface{}
	requestData interface{}
	headers     map[string]string
	HandleFunc  HandleFunc
}

const (
	headerContentType = "Content-Type"
)

func NewHttp() HttpInter {
	return &baseHttp{
		headers: make(map[string]string),
		timeOut: 0,
		HandleFunc: func(ctx context.Context, logger *zap.SugaredLogger, timeOut int64, method string, url string, queryArgs interface{}, requestData interface{}, headers map[string]string) ([]byte, error) {
			var (
				err error
				req *http.Request
			)
			switch requestData.(type) {
			case string:
				req, err = http.NewRequest(method, url, strings.NewReader(requestData.(string)))
			case []byte:
				req, err = http.NewRequest(method, url, bytes.NewBuffer(requestData.([]byte)))
			default:
				req, err = http.NewRequest(method, url, bytes.NewBuffer(nil))
			}
			if err != nil {
				return nil, err
			}
			// post 默认使用 "application/json" 方式
			if method == http.MethodPost {
				if _, ok := headers[headerContentType]; !ok {
					headers[headerContentType] = "application/json"
				}
			}
			if len(headers) > 0 {
				for k := range headers {
					req.Header.Set(k, headers[k])
				}
			}
			req = req.WithContext(ctx)
			if timeOut > 0 {
				childCtx, cancel := context.WithTimeout(ctx, time.Duration(timeOut)*time.Second)
				defer cancel()
				req = req.WithContext(childCtx)
			} else {
				//设置 上下文
				req = req.WithContext(ctx)
			}
			httpTime := time.Now()
			//发起请求的时间
			logger.Infof("[Http] startTime:%d", httpTime.Unix())
			res, err := httpClient.Do(req)
			if err != nil {
				logger.Errorf("[Http] request %+v error:%+v", *req, err)
				return nil, err
			}
			latency := time.Now().Sub(httpTime).Seconds()
			logger.Infof("[Http] latency time consuming:%v", latency)
			defer res.Body.Close()
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logger.Errorf("[Http] read body fail err:%v", err)
				return nil, err
			}
			return b, nil
		},
	}
}
