package report

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"gitlab.yeahka.com/gaas/pkg/util"

	"github.com/ThinkingDataAnalytics/go-sdk/thinkingdata"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/mitchellh/mapstructure"
)

var (
	clients sync.Map
)

type ModeType int32

const (
	ModePatch ModeType = 0 //定时批量写入
	ModeEach  ModeType = 1 //逐条写入
	ModeDebug ModeType = 2 //不入库，只校验数据
)

type Option func(*options)

type options struct {
	mode      ModeType //运行
	eventCfgs map[EventType]EventConfig
	userCfgs  map[UserPropertyType]UserPropertyConfig
}

func WithMode(mode ModeType) Option {
	return func(o *options) {
		o.mode = mode
	}
}

func WithEventCfgs(eventCfgs map[EventType]EventConfig) Option {
	return func(o *options) {
		o.eventCfgs = eventCfgs
	}
}

func WithUserCfgs(userCfgs map[UserPropertyType]UserPropertyConfig) Option {
	return func(o *options) {
		o.userCfgs = userCfgs
	}
}

type thinkDataReporter struct {
	ta        thinkingdata.TDAnalytics
	eventCfgs map[EventType]EventConfig
	userCfgs  map[UserPropertyType]UserPropertyConfig
}

func NewThinkingDataClient(url, appid string, opts ...Option) (Report, error) {
	client, ok := clients.Load(appid)
	if ok {
		if ri, ok := client.(Report); ok {
			return ri, nil
		}
	}
	options := options{
		mode:      ModePatch,
		eventCfgs: map[EventType]EventConfig{},
		userCfgs:  map[UserPropertyType]UserPropertyConfig{},
	}
	for _, o := range opts {
		o(&options)
	}
	var (
		consumer thinkingdata.Consumer
		err      error
	)
	switch options.mode {
	case ModeEach: //逐条传输到TA服务器
		consumer, err = thinkingdata.NewDebugConsumer(url, appid)
	case ModeDebug: //不入库，只校验数据
		consumer, err = thinkingdata.NewDebugConsumerWithWriter(url, appid, false)
	default: //定期定量批量上传
		consumer, err = thinkingdata.NewBatchConsumerWithConfig(thinkingdata.BatchConfig{
			ServerUrl: url,
			AppId:     appid,
			BatchSize: 10,
			AutoFlush: true,
		})
	}
	if err != nil {
		return nil, err
	}
	reporter := &thinkDataReporter{
		ta:        thinkingdata.New(consumer),
		eventCfgs: options.eventCfgs,
		userCfgs:  options.userCfgs,
	}
	clients.Store(appid, reporter)
	return reporter, nil
}

//事件上报
func (t *thinkDataReporter) ReportEvent(ctx context.Context, uxid string, event EventType, data interface{}) error {
	if !t.hasEventConfig(event) {
		return fmt.Errorf("not event cfg for %s", event)
	}
	if !t.isEventData(event, data) {
		return fmt.Errorf("data:%#v not math event:%s cfg", data, event)
	}
	var md map[string]interface{}
	err := mapstructure.Decode(data, &md)
	if err != nil {
		log.Errorf("mapstructure Decode data:%#v fail:%s", data, err)
		return err
	}
	log.Infof("event:%s md:%#v", event, md)
	err = t.ta.Track(uxid, "", string(event), md)
	if err != nil {
		log.Errorf("Track event:%s fail:%s", event, err)
		return err
	}
	return nil
}

//用户属性上报
func (t *thinkDataReporter) ReportUser(ctx context.Context, uxid string, distinctId string, user interface{}) error {
	var md map[string]interface{}
	err := mapstructure.Decode(user, &md)
	if err != nil {
		log.Errorf("Decode user:%#v fail:%s", user, err)
		return err
	}
	onceMap := make(map[string]interface{}) //只设置一次的属性
	setMap := make(map[string]interface{})  //覆盖的属性
	addMap := make(map[string]interface{})  //累加的属性
	for k, v := range md {
		if util.IsBlank(v) {
			continue
		}
		pt := t.userCfgs[UserPropertyType(k)]
		switch pt {
		case UserPropertyTypeOnce:
			onceMap[string(k)] = v
		case UserPropertyTypeSet:
			setMap[string(k)] = v
		case UserPropertyTypeAdd:
			addMap[string(k)] = v
		}
	}
	log.Infof("onceMap:%#v setMap:%#v addMap:%#v", onceMap, setMap, addMap)
	if len(onceMap) > 0 {
		err = t.ta.UserSetOnce(uxid, distinctId, onceMap)
		if err != nil {
			log.Errorf("Track fail:%s", err)
			return err
		}
	}
	if len(setMap) > 0 {
		err = t.ta.UserSet(uxid, distinctId, setMap)
		if err != nil {
			log.Errorf("UserSet fail:%s", err)
			return err
		}
	}
	if len(addMap) > 0 {
		err = t.ta.UserAdd(uxid, distinctId, addMap)
		if err != nil {
			log.Errorf("UserAdd fail:%s", err)
			return err
		}
	}
	return nil
}

//是否做了事件数据类型校验配置
func (t *thinkDataReporter) hasEventConfig(event EventType) bool {
	_, ok := t.eventCfgs[event]
	return ok
}

//判断上报的数据是否为事件所需数据
func (t *thinkDataReporter) isEventData(event EventType, data interface{}) bool {
	if dataTypeCfg, ok := t.eventCfgs[event]; !ok {
		return false
	} else {
		dt := reflect.TypeOf(data)
		cfgt := reflect.TypeOf(dataTypeCfg)
		return dt == cfgt
	}
}
