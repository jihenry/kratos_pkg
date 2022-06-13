package registry

import (
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewNacosClient(addr string, port uint64, copts ...constant.ClientOption) (*nacos.Registry, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(addr, port),
	}
	cc := constant.NewClientConfig(copts...)
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{ClientConfig: cc, ServerConfigs: sc},
	)
	if err != nil {
		return nil, err
	}
	return nacos.New(client), nil
}
