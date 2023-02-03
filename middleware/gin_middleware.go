package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"gitlab.yeahka.com/gaas/pkg/util"

	"gitlab.yeahka.com/gaas/pkg/log"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				data, _ := c.GetRawData()
				buf := make([]byte, 64<<10)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				log.FromContext(c).Errorf("%v: %+v\n%s\n", err, string(data), buf)
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				c.AbortWithStatus(200)
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": http.StatusBadGateway,
					"msg":  http.StatusText(http.StatusBadGateway),
					"data": nil,
				})
				return
			}
		}()
		c.Next()
	}
}

func RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := UUId()
		//日志的上下文
		ctx := c.Request.Context()
		//设置用户的请求ID
		c.Request = c.Request.WithContext(
			util.NewRequestContext(ctx,
				util.RequestData{"requestId": requestId},
			),
		)
		c.Request = c.Request.WithContext(log.NewLoggerContext(c.Request.Context(),
			log.LoggerWith(log.FromContext(ctx), []interface{}{
				RequestID, requestId,
			}...)))
		c.Next()
	}
}

var (
	reqHeaderDump = map[string]bool{
		"Host":              true,
		"Transfer-Encoding": false,
		"Trailer":           false,
		"Accept":            false,
		"Accept-Encoding":   false,
		"Connection":        false,
		"Cache-Control":     true,
		"Accept-Language":   true,
		"Origin":            true,
		"Sec-Fetch-Site":    false,
	}
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger(debug bool, maxByte int, serverName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		data, err := c.GetRawData()
		if err != nil {
			log.FromContext(c).Infof("Logger.GetRawData err:%v", err)
		}
		rw := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = rw
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		c.Next()

		headers := make(http.Header, 0)
		for k, v := range c.Request.Header {
			if h, ok := reqHeaderDump[k]; !(h && ok) {
				headers[k] = v
			}
		}
		var (
			respBody string
		)
		if debug || len(rw.body.Bytes()) <= maxByte {
			respBody = rw.body.String()
		}
		log.FromContext(c.Request.Context()).Infof("[%s] ip:%s | method:%s | uri:%s | headers:%s | req:%s | status:%d | resp:%s | latency:%f",
			serverName,
			c.ClientIP(),
			c.Request.Method,
			c.Request.RequestURI,
			headers,
			string(data),
			c.Writer.Status(),
			respBody,
			time.Since(start).Seconds(),
		)
	}
}
