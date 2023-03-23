package controller

import (
	"bluebell_blogs/logic"
	"bluebell_blogs/models/common"
	"bluebell_blogs/models/dto"
	"bluebell_blogs/models/entity"
	"bluebell_blogs/pkg/getPageInfo"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CommunityHandler 获取社区列表
func CommunityHandler(c *gin.Context) {
	fmt.Println("开始获取社区列表！！！")
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		fmt.Println("获取社区列表数据库操作失败", err)
		common.ResponseError(c, common.CodeServerBusy) //不轻易吧服务端报错暴露给外面
		return
	}
	common.ResponseSuccess(c, data)
}

// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Println("获取社区ID失败 strconv.ParseInt", err)
		common.ResponseError(c, common.CodeInvalidParams)
		return
	}
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		fmt.Println("获取社区分类详情数据库操作失败", err)
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		common.ResponseError(c, common.CodeServerBusy)
		return
	}
	common.ResponseSuccess(c, data)
}

func CreatePostHandler(c *gin.Context) {
	p := new(entity.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) error", zap.Any("err", err))
		zap.L().Error("CreatePost with invalid param")
		fmt.Println("数据接收错误", err)
		common.ResponseError(c, common.CodeInvalidParams)
		return
	}
	//从 C 渠道当前发请求的用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		fmt.Println("获取用户ID错误", err)
		common.ResponseError(c, common.CodeNotLogin)
		return
	}
	p.AuthorID = userID
	//创建帖子
	if err := logic.CreatrPost(p); err != nil {
		fmt.Println("创建帖子错误", err)
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		common.ResponseError(c, common.CodeServerBusy)
		return
	}
	//返回响应
	common.ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {
	//1、获取参数（从URL中获取帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	fmt.Println("pid:", pid)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		fmt.Println(" GetPostDetailHandler", err)
		common.ResponseError(c, common.CodeInvalidParams)
		return
	}
	//2、根据id取出帖子的数据（查数据库）
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed", zap.Error(err))
		fmt.Println("logic.GetPostByID(pid)", err)
		common.ResponseError(c, common.CodeServerBusy)
		return
	}
	//3、返回响应
	common.ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {
	page, size := getPageInfo.GetPage(c)
	fmt.Println("分页参数", page, size)
	//1、获取参数
	data, err := logic.GetPostList(page, size)
	if err != nil {
		fmt.Println("分页错误", err)
		zap.L().Error("logic.GetPostListHandler() failed", zap.Error(err))
		common.ResponseError(c, common.CodeServerBusy)
		return
	}
	common.ResponseSuccess(c, data)
}

// GetPostListHandler2 帖子列表接口
// 根据前端穿传过来的参数动态的获取帖子列表
// 按照创建 时间 排序或者按照 分数 排序
func GetPostListHandler2(c *gin.Context) {
	//1、获取参数
	//2、去redis查询ID列表
	//3、根据ID去数据库查询帖子详情信息
	//获取分页参数
	p := &dto.ParamPostList{
		Page:  1,
		Size:  10,
		Order: dto.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		common.ResponseError(c, common.CodeInvalidParams)
		return
	}
	//c.ShouldBind() 根据请求中的数据类型选择相应的方法去获取数据
	//c.ShouldBingJson() 如果请求中携带的是json格式的数据，才用这个方法获取到数据
	//获取数据
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		common.ResponseError(c, common.CodeServerBusy)
		return
	}
	common.ResponseSuccess(c, data)
}
