package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Mysql struct { // 数据库配置
		DataSource string // mysql链接地址
	}

	CacheRedis cache.CacheConf // redis缓存
}
