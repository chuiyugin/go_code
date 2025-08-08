package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// 声明一个全局rdb变量
var rdb *redis.Client
var ctx = context.Background()

// 初始化连接
func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"), // 密码
		DB:       viper.GetInt("redis.db"),          // 数据库
		PoolSize: viper.GetInt("redis.pool_size"),   // 连接池大小
	})
	_, err = rdb.Ping(ctx).Result()
	return
}

func Close() {
	_ = rdb.Close()
}
