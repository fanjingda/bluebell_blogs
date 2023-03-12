package controller

import (
	"bluebell_blogs/models/common"
	"bluebell_blogs/pkg/jwt"
	"bluebell_blogs/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	user := r.Group("/user")
	{
		user.GET("/hello", func(c *gin.Context) {
			common.ResponseSuccess(c, "ok")
		})

		user.POST("/signup", SignUpHandler)
		user.POST("/login", LoginHandler)
		user.POST("/pang", jwt.JWTAuthMiddleware(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "成功",
			})
		})
	}
	community := r.Group("/community/")
	community.Use(jwt.JWTAuthMiddleware())
	{
		community.POST("/list", CommunityHandler)
		community.POST("/list/:id", CommunityDetailHandler)
		community.POST("/post", CreatePostHandler)
		community.GET("/post/:id", GetPostDetailHandler)
		//根据时间或分数获取帖子列表
		community.GET("/posts", GetPostListHandler)
		community.GET("/posts2", GetPostListHandler2)

	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
