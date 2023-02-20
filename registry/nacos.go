package registry

import (
	"os"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/file"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewNacosClient(addr string, port uint64, copts ...constant.ClientOption) (*nacos.Registry, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(addr, port),
	}
	cc := &constant.ClientConfig{
		TimeoutMs:            10 * 1000,
		BeatInterval:         5 * 1000,
		OpenKMS:              false,
		CacheDir:             file.GetCurrentPath() + string(os.PathSeparator) + "cache",
		UpdateThreadNum:      20,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: false,
		LogDir:               file.GetCurrentPath() + string(os.PathSeparator) + "log",
		RotateTime:           "1h",
		MaxAge:               3,
		LogLevel:             "error",
	}
	for _, opt := range copts {
		opt(cc)
	}
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{ClientConfig: cc, ServerConfigs: sc},
	)
	if err != nil {
		return nil, err
	}
	return nacos.New(client), nil
}
