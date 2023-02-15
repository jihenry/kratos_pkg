## 一、工程规范
1. 存放和服务无关的非业务类代码，包括各类组件接入实例创建、通用算法实现、通用业务功能抽象
2. 此工程基于kratos框架，对于kratos相关子包可以使用，其他框架包禁止直接引入
3. 一级目录用于常用组件，二级目录用于非常用或者是同类组件聚合
## 二、目录说明
```
.
├── README.md
├── cache //存放所有缓存相关组件实例创建
│   ├── lru.go //本地缓存
│   └── redis.go //redis缓存
├── common //存放通用业务功能抽象
│   ├── draw //抽奖
│   └── page //分页
├── db //持续化存储相关组件实例创建
│   ├── log.go
│   └── mysql.go
├── govern //治理相关组件实例创建
│   └── trace.go //链路追踪
├── http //http实例
├── kafka //mq相关组件，目前使用的是kafka
├── log //日志组件实例创建
├── middleware //业务中间件
├── monitor //监控组件相关实例创建
├── party //第三方组件接入
│   ├── cos //腾讯云cos桶
│   ├── openapi //渠道第三方api接入
│   ├── report //上报，目前用的是数数
│   └── weixin //微信
├── registry //注册中心
├── rpc //rpc相关实例
│   └── grpc.go
└── util //通用算法实现
```