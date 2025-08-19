package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// ZRevRange 按照分数从大到小的顺序查询指定数量的元素
	return rdb.ZRevRange(ctx, key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 根据用户请求中携带的
	key := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZset)
	}
	// 确定查询的索引起始点
	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据 ids 查询每篇帖子投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// data = make([]int64, 0, len(ids))
	// for _, id := range ids {
	// 	key := getRedisKey(KeyPostVotedZSetPF + id)
	// 	// 查找key中分数是1的元素的数量 -> 统计每篇帖子的赞成票的数量
	// 	v := rdb.ZCount(ctx, key, "1", "1").Val()
	// 	data = append(data, v)
	// }

	// 使用 pipeline 一次发送多条命令，减少 RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 按社区查询 ids
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZset)
	}
	// 使用 zinterstore 把分区的帖子set与帖子分数的zset生成一个新的zset
	// 针对新的zset 按照之前的逻辑取数据

	// 社区的key
	ckey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存key减少 zinterstore 的执行次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(ctx, key).Val() < 1 {
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Keys:      []string{ckey, orderKey},
			Aggregate: "MAX",
		}) // zinterstore 计算
		pipeline.Expire(ctx, key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	return getIDsFromKey(key, p.Page, p.Size)
}
