package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
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

	// 新增：专门给 token/session 使用的 Redis
	Redis redis.RedisConf

	Auth struct {
		AccessSecret string
		AccessExpire int64
		// 新增：刷新令牌配置
		RefreshSecret string
		RefreshExpire int64
	}
}
