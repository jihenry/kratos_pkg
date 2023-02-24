package monitor

import (
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	srvName, srvip string
	queue          = NewQueue()
	reqMap         sync.Map
	sumMap         sync.Map
	started        = false
	logFile        os.File
)

/*
*
@description：监控日志启动
@param localIp 服务器IP
@param sericeName 服务名称
@param logFilePath 日志文件地址
*/
func Start(localIp string, serviceName string, logFilePath string) {
	if started {
		return
	}
	srvip = localIp
	srvName = serviceName
	started = true
	//初始化日志对象信息
	initLog(serviceName, logFilePath)
	go collectTask()
	go func() {
		ticker := time.NewTicker(2000 * time.Millisecond)
		for started {
			t := <-ticker.C
			outputTask(t)
		}
		ticker.Stop()
	}()
}

/*
*
@description：监控日志停止
*/
func Stop() {
	started = false
	destryLog()
}

type record struct {
	key   string
	vtype uint8
	val   uint64
}

/**
 * @description 记录被调用信息日志。可用于以下场景：
 *   a. 被调用次数统计
 *   b. 被调用耗时统计
 *   c. 被调用错误码(业务响应码)统计
 *   示例场景：HTTP请求、数据库调用记录（可根据具体情况进行自定义扩展customized）、...
 * @param cmd  请求地址/请求指令
 * @param result  请求响应业务错误码或服务错误码（如：0-成功 ...）
 * @param cost  请求耗时情况（ms）
 * @param customized  自定义键值日志，可根据具体场景定义（建议：键值不宜过多、同一CMD下尽量统一规范）
 */
func Req(cmd string, result string, cost uint64, customized map[string]string) {
	var reqstr string
	var sumstr string
	if !started {
		return
	}
	if customized != nil && len(customized) > 0 {
		var tstr string
		for item, val := range customized {
			tstr += "," + item + "=\"" + val + "\""
		}
		reqstr = fmt.Sprintf("ykgame_req{cmd=\"%v\",ipp=\"%v\",errcode=\"%v\",srv=\"%v\"%v}", cmd, srvip, result, srvName, tstr)
		sumstr = fmt.Sprintf("ykgame_req_summary{cmd=\"%v\",ipp=\"%v\",errcode=\"%v\",srv=\"%v\"%v,percentile=\"%%v\"}", cmd, srvip, result, srvName, tstr)
	} else {
		reqstr = fmt.Sprintf("ykgame_req{cmd=\"%v\",ipp=\"%v\",errcode=\"%v\",srv=\"%v\"}", cmd, srvip, result, srvName)
		sumstr = fmt.Sprintf("ykgame_req_summary{cmd=\"%v\",ipp=\"%v\",errcode=\"%v\",srv=\"%v\",percentile=\"%%v\"}", cmd, srvip, result, srvName)
	}
	recordReq := record{reqstr, 1, 1}
	recordSum := record{sumstr, 2, cost}
	queue.PushBack(recordReq)
	queue.PushBack(recordSum)
}

/**
 * @description  记录远程调用资源信息日志。可用于以下场景：
 *   a.远程调用第三方接口次数统计
 *   b.远程调用第三方接口耗时统计
 *   b.远程调用第三方接口响应状态统计
 *   ...
 * @param cmd  第三方调用指令/请求地址/方法名称（例：hsycouponweb.yeahka.com/xfy/getMerchanid）
 * @param result  请求响应业务错误码或服务错误码（如：0-成功 ...）
 * @param cost  请求耗时情况（ms）
 * @param rmtsrv  远程服务名称（例如:hsycouponweb.yeahka.com或hsycouponweb.yeahka.com/xfy）
 * @param rmtipp  远程服务IP（例：:hsycouponweb.yeahka.com）
 */
func Rpc(cmd string, result string, cost uint64, rmtsrv string, rmtipp string) {
	if !started {
		return
	}
	reqstr := fmt.Sprintf("ykgame_rpc{cmd=\"%v\",ipp=\"%v\",errcode=\"%v\",srv=\"%v\",rmtsrv=\"%v\",rmtipp=\"%v\"}", cmd, srvip, result, srvName, rmtsrv, rmtipp)
	sumstr := fmt.Sprintf("ykgame_rpc_summary{cmd=\"%v\",ipp=\"%v\",errcode=\"%v\",srv=\"%v\",rmtsrv=\"%v\",rmtipp=\"%v\",percentile=\"%%v\"}", cmd, srvip, result, srvName, rmtsrv, rmtipp)
	recordReq := record{reqstr, 1, 1}
	recordSum := record{sumstr, 2, cost}
	queue.PushBack(recordReq)
	queue.PushBack(recordSum)
}

/**
 * @description 记录资源使用（消耗）情况日志。可用于以下场景：
 *   a.数据库连接池使用情况
 *   b.Redis连接池使用情况
 *   c.FASTDFS连接池使用情况
 *   ...
 * @param src  资源名称，如： mysql_pool_usage
 * @param identity  具体资源对象，如：连接池名称
 * @param used	资源已使用量
 * @param maxSrc  资源最大量
 */
func Src(src string, identity string, used uint64, maxSrc uint64) {
	if !started {
		return
	}
	sumstr := fmt.Sprintf("ykgame_src{ipp=\"%v\",maxSrc=\"%v\",percentile=\"%%v\",process=\"%v\",src=\"%v\",srv=\"%v\"}", srvip, maxSrc, identity, src, srvName)
	maxSumstr := fmt.Sprintf("ykgame_src_max{ipp=\"%v\",maxSrc=\"%v\",percentile=\"%%v\",process=\"%v\",src=\"%v\",srv=\"%v\"}", srvip, maxSrc, identity, src, srvName)
	recordSum := record{sumstr, 2, used}
	recordMaxSum := record{maxSumstr, 2, maxSrc}
	queue.PushBack(recordSum)
	queue.PushBack(recordMaxSum)
}

// 收集任务
func collectTask() {
	for started {
		item := queue.Front()
		if item == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		record := (item.Value).(record)
		if record.vtype == 1 {
			value, _ := reqMap.Load(record.key)
			if value != nil {
				reqMap.Store(record.key, value.(uint64)+record.val)
			} else {
				reqMap.Store(record.key, record.val)
			}

		} else if record.vtype == 2 {
			value, _ := sumMap.Load(record.key)
			if value != nil {
				sumMap.Store(record.key, append(value.([]uint64), record.val))
			} else {
				sumMap.Store(record.key, []uint64{record.val})
			}
		}
		queue.Remove(item)
	}
}

// 输出任务
func outputTask(t time.Time) {
	dstr := t.Format("2006-01-02 15:04:05")
	rstr := dstr + "\n"
	count := 0
	reqMap.Range(func(key, value interface{}) bool {
		rstr += fmt.Sprintf("%v %v\n", value.(uint64), key)
		count++
		return true
	})
	if count > 0 {
		log.Info(rstr)
	}
	reqMap = sync.Map{}

	sstr := dstr + "\n"
	count = 0
	sumMap.Range(func(key, value interface{}) bool {
		sstr += fmt.Sprintf("%v "+key.(string)+"\n", calcQuantile(value.([]uint64)), "P99")
		count++
		return true
	})
	if count > 0 {
		log.Info(sstr)
	}
	sumMap = sync.Map{}
}

// P99中位数计算
func calcQuantile(darr []uint64) uint64 {
	arrlen := len(darr)
	if arrlen == 0 {
		return 0
	}
	if arrlen == 1 {
		return darr[0]
	}
	sort.Slice(darr, func(x, y int) bool {
		return darr[x] < darr[y]
	})
	a := float64(arrlen-1) * 0.99
	n := math.Floor(a)
	f := a - n
	result := (1-f)*float64(darr[uint32(n)]) + f*float64(darr[uint32(n)+1])
	return uint64(math.Floor(result + 0.5))
}

func destryLog() {
	logFile.Close()
}
