package mysql

import (
	"bluebell_blogs/models/common"
	"bluebell_blogs/models/entity"
	"database/sql"
	"go.uber.org/zap"
)

func GetCommunityList() (CommunityList []*entity.Community, err error) {
	sqlStr := `select community_id,community_name from community`
	if err := db.Select(&CommunityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

// GetCommunityDetailByID 根据ID查询分类社区详情
func GetCommunityDetailByID(id int64) (community *entity.CommunityDetail, err error) {
	community = new(entity.CommunityDetail)
	sqlStr := `select community_id,community_name,introduction,create_time from community where community_id=?`
	if err := db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = common.ErrorInvalidID
		}
	}
	return community, err
}
