package logic

import (
	"bluebell_blogs/dao/mysql"
	"bluebell_blogs/models/entity"
)

func GetCommunityList() ([]*entity.Community, error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*entity.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
