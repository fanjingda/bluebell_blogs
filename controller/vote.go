package controller

import (
	"bluebell_blogs/logic"
	"bluebell_blogs/models/common"
	"bluebell_blogs/models/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

//投票

func PostVoteController(c *gin.Context) {
	//参数校验
	p := new(dto.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		fmt.Println("获取参数错误", err)
		zap.L().Error("投票参数获取失败", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			common.ResponseErrorWithMsg(c, common.CodeInvalidParams, errs)
			return
		}
		errData := RemoveTopStruct(errs.Translate(Trans))
		common.ResponseErrorWithMsg(c, common.CodeInvalidParams, errData)
		return
	}
	userID, err := getCurrentUserID(c)
	if err != nil {
		common.ResponseError(c, common.CodeNotLogin)
		return
	}
	if err := logic.PostVote(userID, p); err != nil {
		fmt.Println("投票数据库错误", err)
		zap.L().Error("投票业务逻辑错误", zap.Error(err))
		common.ResponseErrorWithMsg(c, common.CodeServerBusy, err)
		return
	}
	common.ResponseSuccess(c, nil)
}
