package nacos

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/spf13/viper"

	"gitlab.yeahka.com/gaas/pkg/zaplog"

	"google.golang.org/grpc"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	mmd "github.com/go-kratos/kratos/v2/middleware/metadata"
	tgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type (
	NacosConfig struct {
		NamespaceID         string       `yaml:"namespaceId"`
		TimeoutMs           uint64       `yaml:"timeoutMs"`
		NotLoadCacheAtStart bool         `yaml:"notLoadCacheAtStart"`
		LogDir              string       `yaml:"logDir"`
		CacheDir            string       `yaml:"cacheDir"`
		LogLevel            string       `yaml:"logLevel"`
		RotateTime          string       `yaml:"rotateTime"`
		MaxAge              int64        `yaml:"maxAge"`
		NacosServer         ServerConfig `yaml:"nacosServer"`
	}
	ServerConfig struct {
		IP   string `yaml:"ip"`
		Port uint64 `yaml:"port"`
	}
)

func getClientConfig(cc NacosConfig) (*constant.ClientConfig, []constant.ServerConfig) {
	return &constant.ClientConfig{
			NamespaceId:         cc.NamespaceID,
			TimeoutMs:           cc.TimeoutMs,
			NotLoadCacheAtStart: cc.NotLoadCacheAtStart,
			LogDir:              cc.LogDir,
			CacheDir:            cc.CacheDir,
			LogLevel:            cc.LogLevel,
			RotateTime:          cc.RotateTime,
			MaxAge:              cc.MaxAge,
		}, []constant.ServerConfig{
			*constant.NewServerConfig(cc.NacosServer.IP, cc.NacosServer.Port),
		}
}

var (
	nacClient           *nacos.Registry
	nacMutex            sync.Mutex
	nacosConfig         map[string]*viper.Viper
	nacConfig           config_client.IConfigClient
	nacConfigMutex      sync.Mutex
	naccosConfigRwMutex sync.RWMutex
	commonConfigRwMutex sync.RWMutex
)

func init() {
	nacosConfig = make(map[string]*viper.Viper, 0)
}

type ConfigParam struct {
	DataId string `yaml:"dataId"` //required
	Group  string `yaml:"group"`  //required
	Name   string `yaml:"name"`
}

type NacosListenConfig struct {
	DataId   string
	Group    string
	OnChange func() func(namespace, group, dataId, data string)
}

func GetConfig(name string) *viper.Viper {
	defer naccosConfigRwMutex.RUnlock()
	naccosConfigRwMutex.RLock()
	conf, ok := nacosConfig[name]
	if !ok || conf == nil {
		panic(errors.New(fmt.Sprintf("config is nil")))
	}
	return conf
}

func newConfigClient(nc NacosConfig) (config_client.IConfigClient, error) {
	if nacConfig == nil {
		nacConfigMutex.Lock()
		defer nacConfigMutex.Unlock()
		if nacConfig == nil {
			cc, sc := getClientConfig(nc)
			//服务注册, 服务连接
			cli, err := clients.NewConfigClient(vo.NacosClientParam{
				ClientConfig:  cc,
				ServerConfigs: sc,
			})
			if err != nil {
				return nil, err
			}
			nacConfig = cli
		}
	}
	return nacConfig, nil
}

func NewEnvConfig(nc NacosConfig, dataId, group string, fun func(str string) error, reLoadConfig func(str string) error) error {
	ncl, err := newConfigClient(nc)
	if err != nil {
		return err
	}
	content, err := ncl.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return err
	}
	if err := fun(content); err != nil {
		return err
	}
	go func() {
		_ = ncl.ListenConfig(vo.ConfigParam{
			DataId: dataId,
			Group:  group,
			OnChange: func(namespace, group, dataId, data string) {
				if err := reLoadConfig(data); err != nil {
					log.Printf("err:%v", err)
				}
			},
		})
	}()
	return nil
}

func NewConfigClient(nc NacosConfig, cp []ConfigParam) error {
	ctx := context.Background()
	ncl, err := newConfigClient(nc)
	if err != nil {
		return err
	}
	for k := range cp {
		defaultConfig := viper.New()
		defaultConfig.SetConfigType("json")
		param := cp[k]
		content, err := ncl.GetConfig(vo.ConfigParam{
			DataId: param.DataId,
			Group:  param.Group,
		})
		if err != nil {
			return err
		}
		if err := defaultConfig.ReadConfig(bytes.NewReader([]byte(content))); err != nil {
			return err
		}
		nacosConfig[param.Name] = defaultConfig
	}
	go func() {
		for k := range cp {
			defaultConfig := viper.New()
			param := cp[k]
			defaultConfig.SetConfigType("json")
			if err := ncl.ListenConfig(vo.ConfigParam{
				DataId: param.DataId,
				Group:  param.Group,
				OnChange: func(namespace, group, dataId, data string) {
					if err := defaultConfig.ReadConfig(bytes.NewReader([]byte(data))); err != nil {
						zaplog.FromContext(ctx).Errorf("readConfig fail err:%v", err)
					}
					naccosConfigRwMutex.Lock()
					nacosConfig[param.Name] = defaultConfig
					naccosConfigRwMutex.Unlock()
				},
			}); err != nil {
				zaplog.FromContext(ctx).Errorf("biper fail %v", err)
			}
		}
	}()
	return nil
}

// RegistryNacos 服务注册与发现
func RegistryNacos(nc NacosConfig) (*nacos.Registry, error) {
	if nacClient == nil {
		nacMutex.Lock()
		defer nacMutex.Unlock()
		if nacClient == nil {
			cc, sc := getClientConfig(nc)
			//服务注册, 服务连接
			cli, err := clients.NewNamingClient(
				vo.NacosClientParam{
					ClientConfig:  cc,
					ServerConfigs: sc,
				},
			)
			if err != nil {
				return nil, err
			}
			nacClient = nacos.New(cli)
		}
	}
	return nacClient, nil
}

type GRpcConn interface {
	Close()
}

var (
	_         GRpcConn = (*conn)(nil)
	clientMap          = map[string]GRpcConn{}
	connMap            = map[string]*grpc.ClientConn{}
)

func init() {
	clientMap = make(map[string]GRpcConn, 0)
	connMap = make(map[string]*grpc.ClientConn, 0)
}

type conn struct {
	sync.Mutex
}

func newGRpcConn(name string, gc *grpc.ClientConn) (GRpcConn, error) {
	c := new(conn)
	c.setConn(name, gc)
	return c, nil
}

//GetConn  获取链接
func GetConn(name string) *grpc.ClientConn {
	return connMap[name]
}

//设置链接
func (c *conn) setConn(name string, conn *grpc.ClientConn) {
	c.Lock()
	defer c.Unlock()
	connMap[name] = conn
}

func (c *conn) Close() {
	for k := range connMap {
		connMap[k].Close()
	}
}

func clientConn(ctx context.Context, endpoint string, nc NacosConfig) (GRpcConn, error) {
	var (
		nr  *nacos.Registry
		err error
		cn  *grpc.ClientConn
	)
	if nr, err = RegistryNacos(nc); err != nil {
		return nil, err
	}
	if cn, err = tgrpc.DialInsecure(
		ctx,
		//endpoint
		tgrpc.WithEndpoint(fmt.Sprintf("discovery:///%s.grpc", endpoint)),
		//注册服务
		tgrpc.WithDiscovery(nr),
		tgrpc.WithMiddleware(
			mmd.Client(),
			mmd.Server(),
			tracing.Client(),
			tracing.Server(),
		),
		tgrpc.WithTimeout(0),
	); err != nil {
		return nil, err
	}
	return newGRpcConn(endpoint, cn)
}

//InitNacosClient  初始化服务
func InitNacosClient(nc NacosConfig, name ...string) (func(), error) {
	ctx := context.Background()
	logger := zaplog.LoggerWith(zaplog.FromContext(ctx), []interface{}{"InitGRPCServer", "InitNacosClient"}...)
	for _, v := range name {
		var (
			gc  GRpcConn
			err error
		)
		if gc, err = clientConn(ctx, v, nc); err != nil {
			logger.Errorf("get client conn fail err:%v", err)
			continue
		}
		clientMap[v] = gc
	}
	return func() {
		for k := range clientMap {
			clientMap[k].Close()
		}
	}, nil
}
