package logic

import (
	"bluebell_blogs/dao/redis"
	"bluebell_blogs/models/dto"
	"go.uber.org/zap"
	"strconv"
)

func PostVote(userID int64, p *dto.ParamVoteData) error {
	zap.L().Debug("PostVote failed", zap.Int64("userID", userID), zap.String("postID", p.PostID))
	return redis.PostForVote(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
