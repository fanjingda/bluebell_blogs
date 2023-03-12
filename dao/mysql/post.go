package mysql

import (
	"bluebell_blogs/models/entity"
	"github.com/jmoiron/sqlx"
	"strings"
)

func CreatPost(p *entity.Post) (err error) {
	sqlStr := `insert into 
    			post(post_id,title,content,author_id,community_id)
				values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostByID 根据ID查询单个帖子数据
func GetPostByID(pid int64) (post *entity.Post, err error) {
	post = new(entity.Post)
	sqlStr := `select post_id,title,content,author_id,community_id
from post where post_id =?`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 查询帖子列表函数
func GetPostList(page, size int64) (posts []*entity.Post, err error) {
	sqlStr := `select 
    post_id,title,content,author_id,community_id,create_time
	from post 
	Order By  create_time
	DESC 
	limit ?,?`
	posts = make([]*entity.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDS 根据给定的ID列表查询帖子数据
func GetPostListByIDS(ids []string) (postList []*entity.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
			from post
			where post_id in (?)
			order by FIND_IN_SET(post_id,?)
			`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...) //!!!!!
	return
}
