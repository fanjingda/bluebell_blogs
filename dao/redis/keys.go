package redis

const (
	KeyPrefix              = "bluebell:"
	KeyPostTimeZSet        = "post:time"   //ZSet;帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  //ZSet;帖子及投票的分数
	KeyPostVotedZSetPrefix = "post:voted:" //ZSet;记录用户及投票类型，参数是post_id
)
