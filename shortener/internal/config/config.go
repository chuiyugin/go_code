package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	ShortUrlDB struct {
		DSN string
	}

	SequenceDB struct {
		DSN string
	}

	ShortUrlBlackList []string

	ShortDoamin string

	CacheRedis cache.CacheConf // redis缓存

	Auth struct { // JWT 认证需要的密钥和过期时间配置
		AccessSecret string
		AccessExpire int64
	}
}
