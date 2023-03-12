package redis

import (
	"bluebell_blogs/models/common"
	"errors"
	"math"
	"time"

	"github.com/go-redis/redis"
)

/*
	direction =1 时，有两种情况:
		1、之前没有投过票，现在投赞成票 --->更新分数和投票记录   差值的绝对值：1
		2、之前投反对票，现在改投赞成票 --->更新分数和投票记录   差值的绝对值：2
	direction =0时 ， 有两种情况
		1、之前投赞成票，现在要取消投票  --->更新分数和投票记录  差值的绝对值：1
		2、之前投过反对票，现在要取消投票 --->更新分数和投票记录  差值的绝对值：1
	direction = -1 时，有两种情况
		1、之前没有投过票，现在投反对票  --->更新分数和投票记录  差值的绝对值：1
		2、之前头皮赞成票，现在改投反对票 --->更新分数和投票记录  差值的绝对值：2

	投票限制:
		每个帖子自发表之日起一个星期之内允许用户图片，超过一个星期就不允许在投票了
		1、到期之后将Redis中保存的赞成票票数以及反对票数存储到mysql表中
		2、到期之后删除那个KeyPostVotedZSetPf
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票多少分
)

var (
	ErrVoteTimeExpire = errors.New("超出投票时间")
	ErrVoteRepested   = errors.New("不允许重复投票")
)

func CreatPost(postID int64) error {
	pipeline := client.TxPipeline()
	pipeline.ZAdd(common.GetRedisKey(common.KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	pipeline.ZAdd(common.GetRedisKey(common.KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	_, err := pipeline.Exec()
	return err
}

func PostForVote(userID, postID string, v float64) error {
	//1、判读投票的限制
	//去Redis取帖子发布时间
	postTime := client.ZScore(common.GetRedisKey(common.KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2、更新分数
	//先查当前用户之前给帖子的投票记录
	ov := client.ZScore(common.GetRedisKey(common.KeyPostVotedZSetPf+postID), userID).Val()
	//果然这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if v == ov {
		return ErrVoteRepested
	}
	var dir float64
	if v > ov {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(ov - v) //计算分数差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(common.GetRedisKey(common.KeyPostScoreZSet),
		dir*diff*scorePerVote,
		postID)

	//3、记录用户未改帖子投票的数据
	if v == 0 {
		pipeline.ZRem(common.GetRedisKey(common.KeyPostVotedZSetPf+postID), userID)
	} else {
		pipeline.ZAdd(common.GetRedisKey(common.KeyPostVotedZSetPf+postID), redis.Z{
			Score:  v,
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
