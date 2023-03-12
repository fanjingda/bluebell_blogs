package redis

import (
	"bluebell_blogs/models/common"
	"bluebell_blogs/models/dto"

	"github.com/go-redis/redis"
)

func GetPostIDInOrder(p *dto.ParamPostList) ([]string, error) {
	//从redis获取id
	//1、根据用户请求中携带的order参数确定要查询的redis key
	key := common.GetRedisKey(common.KeyPostTimeZSet)
	if p.Order == dto.OrderScore {
		key = common.GetRedisKey(common.KeyPostScoreZSet)
	}
	//2、确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	//3、ZRevRange 按分数从大到小查询
	return client.ZRevRange(key, start, end).Result()
}

// GetPostVoteDate 根据ids查询每一篇帖子的投赞成票的数据
func GetPostVoteDate(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPf + id)
	//	v1 := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v1)
	//}
	//使用pipeline一次发送多条命令减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := common.GetRedisKey(common.KeyPostVotedZSetPf + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
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

////func GetCommunityPostIDInOrder(p *models.ParamCommunityPostList) (data []int64, err error) {
//	key := getRedisKey(KeyPostTimeZSet)
//	if p.Order == models.OrderScore {
//		key = getRedisKey(KeyPostScoreZSet)
//	}
//	//2、确定查询的索引起始点
//	start := (p.Page - 1) * p.Size
//	end := start + p.Size - 1
//	//3、ZRevRange 按分数从大到小查询
//	return client.ZRevRange(key, start, end)
//}
