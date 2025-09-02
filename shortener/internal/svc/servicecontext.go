package svc

import (
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	ShortUrlModel model.ShortUrlMapModel

	Sequence sequence.Sequence

	ShortUrlBlackList map[string]struct{}

	// 布隆过滤器
	Filter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	// 把配置文件中的黑名单加载到map，方便后续查询
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}
	// 初始化布隆过滤器
	store := redis.New(c.CacheRedis[0].Host, func(r *redis.Redis) {
		r.Type = redis.NodeType
	})
	// 声明一个 filter, key="test_key"且bits是1024位
	filter := bloom.New(store, "test_key", 1024)

	return &ServiceContext{
		Config:        c,
		ShortUrlModel: model.NewShortUrlMapModel(conn, c.CacheRedis),

		Sequence: sequence.NewMySQL(c.SequenceDB.DSN),

		ShortUrlBlackList: m, // 短链接黑名单map

		Filter: filter,
	}
}
