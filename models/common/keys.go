package common

const (
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "post:time"   //zset;帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  //zset;帖子及投票分数
	KeyPostVotedZSetPf = "post:voted:" //记录用及投票类型；参数是post_id

	KeyCommunitySetPf = "community" //set;保存每个分区下的帖子ID
)

// GetRedisKey 加上前缀
func GetRedisKey(key string) string {
	return Prefix + key
}
