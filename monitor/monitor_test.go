package monitor

import (
	"testing"
	"time"
)

func TestMonitor(t *testing.T) {
	Start("127.0.0.1", "micro", "logs")

	Req("aa", "200", 1000, nil)
	Req("aa", "200", 600, nil)
	Req("aa", "200", 1300, nil)
	Req("bb", "300", 3000, nil)

	Rpc("aa", "200", 1000, "www.baidu.com", "199.99.99.99")
	Rpc("aa", "200", 600, "www.baidu.com", "199.99.99.99")
	Rpc("aa", "200", 1300, "www.baidu.com", "199.99.99.99")
	Rpc("bb", "300", 3000, "www.baidu.com", "199.99.99.99")

	Src("mysql_pool_usage", "Mysql_Pool", 12, 30)

	time.Sleep(6 * time.Second)
	Stop()
}
