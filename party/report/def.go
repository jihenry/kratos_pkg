package report

import "context"

//EventConfig 事件类型
type EventType string

//EventConfig 事件对应的值的元数据配置
type EventConfig interface{}

//UserPropertyType 用户属性类型
type UserPropertyType string

//UserPropertyConfig 用户属性元数据配置
type UserPropertyConfig int8

const (
	UserPropertyTypeOnce UserPropertyConfig = iota //第一次属性
	UserPropertyTypeSet                            //覆盖属性
	UserPropertyTypeAdd                            //累加属性
)

type Report interface {
	ReportEvent(ctx context.Context, uxid string, event EventType, data interface{}) error
	ReportUser(ctx context.Context, uxid string, distinctId string, user interface{}) error
}
