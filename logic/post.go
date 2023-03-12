package logic

import (
	"bluebell_blogs/dao/mysql"
	"bluebell_blogs/dao/redis"
	"bluebell_blogs/models/dto"
	"bluebell_blogs/models/entity"
	"bluebell_blogs/pkg/snowflake"
	"fmt"
	"go.uber.org/zap"
)

func CreatrPost(p *entity.Post) (err error) {
	//1、生成post ID
	p.ID = snowflake.GenID()
	//2、保存到数据库
	err = mysql.CreatPost(p)
	if err != nil {
		fmt.Println("创建帖子数据库错误", err)
		return err
	}
	err = redis.CreatPost(p.ID)
	return
	//3、返回
}
func GetPostByID(pid int64) (data *entity.ApiPostDetail, err error) {
	//查询并组合接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		fmt.Println("通过P.ID查询错误", err)
		zap.L().Error("mysql.GetPostByID(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	//根据作者id查询作者ID
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		fmt.Println("通过作者ID查询错误", err)
		zap.L().Error("mysql.GetPostByID(post.AuthorID) failed", zap.Error(err))
		return
	}
	//根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		fmt.Println("通过社区ID查询错误", err)
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
		return
	}
	data = &entity.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*entity.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*entity.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetPostByID(post.AuthorID) failed", zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
			continue
		}
		postDetail := &entity.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *dto.ParamPostList) (data []*entity.ApiPostDetail, err error) {
	//2、去redis查询ID列表
	//3、根据ID去数据库查询帖子详情信息
	ids, err := redis.GetPostIDInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDInOrder(p) return 0 data")
		return
	}
	posts, err := mysql.GetPostListByIDS(ids)
	if err != nil {
		return
	}
	voteData, err := redis.GetPostVoteDate(ids)
	if err != nil {
		return
	}
	//将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetPostByID(post.AuthorID) failed", zap.Error(err))
			continue
		}
		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
			continue
		}
		postDetail := &entity.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
