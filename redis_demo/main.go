package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 声明一个全局rdb变量
var rdb *redis.Client
var ctx = context.Background()

// 初始化连接
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 密码
		DB:       0,  // 数据库
		PoolSize: 20, // 连接池大小
	})
	_, err = rdb.Ping(ctx).Result()
	return err
}

func redisExample() {
	err := rdb.Set(ctx, "score", 100, 0).Err()
	if err != nil {
		fmt.Printf("set score failed! err:%#v\n", err)
	}

	val, err := rdb.Get(ctx, "score").Result()
	if err != nil {
		fmt.Printf("get score failed! err:%#v\n", err)
	}
	fmt.Println("score:", val)

	val_2, err := rdb.Get(ctx, "name").Result()
	if err == redis.Nil {
		fmt.Println("score does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("score:", val_2)
	}
}

// zsetDemo 操作zset示例
func zsetDemo() {
	// key
	zsetKey := "language_rank"
	// value
	// 注意：v8版本使用[]*redis.Z；此处为v9版本使用[]redis.Z
	languages := []redis.Z{
		{Score: 90.0, Member: "Golang"},
		{Score: 98.0, Member: "Java"},
		{Score: 95.0, Member: "Python"},
		{Score: 97.0, Member: "JavaScript"},
		{Score: 99.0, Member: "C/C++"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// ZADD
	err := rdb.ZAdd(ctx, zsetKey, languages...).Err()
	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}
	fmt.Println("zadd success")

	// 把Golang的分数加10
	newScore, err := rdb.ZIncrBy(ctx, zsetKey, 10.0, "Golang").Result()
	if err != nil {
		fmt.Printf("zincrby failed, err:%v\n", err)
		return
	}
	fmt.Printf("Golang's score is %f now.\n", newScore)

	// 取分数最高的3个
	ret := rdb.ZRevRangeWithScores(ctx, zsetKey, 0, 2).Val()
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}

	// 取95~100分的
	op := &redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}
	ret, err = rdb.ZRangeByScoreWithScores(ctx, zsetKey, op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}
}

func main() {
	if err := initClient(); err != nil {
		fmt.Printf("init redis client failed! err: %#v\n", err)
		return
	}
	fmt.Println("init redis sucess!")
	// 程序退出时释放相关资源
	defer rdb.Close()

	// redisExample()
	zsetDemo()
}
