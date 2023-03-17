package dto

// 路由输入参数
const (
	OrderTime  = "time"
	OrderScore = "score"
)

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"RePassword" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
}

type ParamVoteData struct {
	//UserID从请求中获取当前的用户
	PostID    string `json:"post_id" binding:"required"`                       //帖子ID
	Direction int    `json:"direction,string" binding:"required,oneof=1 0 -1"` //赞成票（1）反对票（-1）
}
