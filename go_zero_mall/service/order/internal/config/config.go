package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserRPC zrpc.RpcClientConf // 连接其他微服务的RPC客户端

	Mysql struct { // 数据库配置
		DataSource string // mysql链接地址
	}

	CacheRedis cache.CacheConf // redis缓存
}
