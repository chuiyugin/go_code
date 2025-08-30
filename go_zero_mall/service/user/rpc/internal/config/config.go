package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct { // 数据库配置
		DataSource string // mysql链接地址
	}

	CacheRedis cache.CacheConf // redis缓存

	Consul consul.Conf
}
