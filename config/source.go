package config

import (
	"fmt"
	"sync"

	nacos "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
)

type ConfigOption func(*configOpt)

type configOpt struct {
	nacosClient config_client.IConfigClient
}

func WithClient(client config_client.IConfigClient) ConfigOption {
	return func(co *configOpt) {
		co.nacosClient = client
	}
}

var (
	configMap sync.Map
)

func NewNacosConfig(dataId, groupId string, opts ...ConfigOption) (config.Config, error) {
	if dataId == "" || groupId == "" {
		return nil, fmt.Errorf("dataId:%s or groupId:%s can't empty", dataId, groupId)
	}
	option := configOpt{
		nacosClient: pvNacosConfigClient,
	}
	for _, opt := range opts {
		opt(&option)
	}
	key := dataId + ":" + groupId
	if configIns, ok := configMap.Load(fmt.Sprintf("%s:%s", dataId, groupId)); ok {
		if ins, ok := configIns.(config.Config); ok {
			return ins, nil
		}
	}
	cfgIns := config.New(
		config.WithSource(nacos.NewConfigSource(option.nacosClient, nacos.WithDataID(dataId), nacos.WithGroup(groupId))),
	)
	if err := cfgIns.Load(); err != nil {
		return nil, err
	}
	configMap.Store(key, cfgIns)
	return cfgIns, nil
}
