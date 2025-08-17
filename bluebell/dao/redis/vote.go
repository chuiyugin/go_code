package redis

import (
	"errors"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

// 本项目采用简化版的投票分数
// 投一票加432分 86400/200 -> 200张赞成票可以给帖子续一天 -> 《redis实战》

/*
投票的几种情况：
direction=1时，有两种情况：
	1、之前没有投过票，现在投赞成票 --> 更新分数和投票记录 差值的绝对值：1 +432
	2、之前投反对票，现在改投赞成票 --> 更新分数和投票记录 差值的绝对值：2 +432*2
direction=0时，有两种情况：
	1、之前投反对票，现在取消投票   --> 更新分数和投票记录 差值的绝对值：1 +432
	2、之前投过赞成票，现在取消投票 --> 更新分数和投票记录 差值的绝对值：1 -432
direction=-1时，有两种情况：
	1、之前没有投过票，现在投反对票 --> 更新分数和投票记录 差值的绝对值：1 -432
	2、之前投赞成票，现在改投反对票 --> 更新分数和投票记录 差值的绝对值：2 -432*2
投票的限制：
每个帖子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1、到期之后redis中保存的赞成票数以及反对票数存储到mysql表中
	2、到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
)

func CreatePost(postID int64) error {
	// 放到一个事务里面，保证同时成功
	pipeline := rdb.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 帖子分数
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipeline.Exec(ctx)
	return err
}

func VoteforPost(userID, postID string, value float64) error {
	// 1. 判断投票限制
	// 去redis取帖子的发布时间
	postTime := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZset), postID).Val()
	if time.Now().Unix()-int64(postTime) > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 2. 更新帖子的分数
	// 先查当前用户给当前帖子的投票记录
	ov := rdb.ZScore(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	var op float64
	var err error
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	// 2和3的操作需要放到一个pipeline事务中操作
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZset), op*diff*scorePerVote, postID)
	// 3. 记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZSetPF+postID), postID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, // 赞成票还是反对票
			Member: userID,
		})
	}
	_, err = pipeline.Exec(ctx)
	return err
}
