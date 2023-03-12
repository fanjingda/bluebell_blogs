package controller

import (
	"bluebell_blogs/logic"
	"bluebell_blogs/models/common"
	"bluebell_blogs/models/dto"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 获取注册相关参数与参数校验
func SignUpHandler(c *gin.Context) {
	//获取参数和参数校验
	p := new(dto.ParamSignUp)
	fmt.Println(" 注册用户开始获取参数和参数校验")
	if err := c.ShouldBindJSON(p); err != nil {
		fmt.Println(err)
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			fmt.Println(errs)
			common.ResponseError(c, common.CodeInvalidParams)
			return
		}
		common.ResponseErrorWithMsg(c, common.CodeInvalidParams, RemoveTopStruct(errs.Translate(Trans)))
		return
	}
	//业务处理
	fmt.Println(" 注册用户业务处理")
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, common.ErrorUserExist) {
			fmt.Println("  mysql.ErrorUserExist;数据库错误", err)
			common.ResponseError(c, common.CodeUserExist)
			return
		}
		fmt.Println("logic.SignUp(p)", err)
		common.ResponseError(c, common.CodeServerBusy)
		return
	}
	//返回响应
	common.ResponseSuccess(c, common.CodeSuccess)
}

// LoginHandler 获取登录相关参数与参数校验
func LoginHandler(c *gin.Context) {
	p := new(dto.ParamLogin)
	fmt.Println("开始登陆操作")
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误
		zap.L().Error("Login with invalid param", zap.Error(err))
		fmt.Println("请求参数有误", err)
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			common.ResponseError(c, common.CodeInvalidParams)
			return
		}
		common.ResponseErrorWithMsg(c, common.CodeInvalidParams, RemoveTopStruct(errs.Translate(Trans)))
		return
	}

	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic Login failed", zap.String("username", p.Username), zap.Error(err))
		fmt.Println("logic Login failed", err)
		if errors.Is(err, common.ErrorInvalidPassword) {
			common.ResponseError(c, common.CodeUserNotExist)
			return
		}
		common.ResponseError(c, common.CodeInvalidPassword)
		return
	}
	fmt.Println(token)
	common.ResponseSuccess(c, token)
}

// getCurrentUserID 通过c.Get来获取保存在token中的userID
func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	_userID, ok := c.Get(common.CUserIDKey)
	if !ok {
		fmt.Println("没有获取到用户ID", ok)
		err = common.ErrorUserNotLogin
		return
	}
	fmt.Println("获取到的用户ID为:", _userID)
	userID, ok = _userID.(int64)
	if !ok {
		fmt.Println("没有用户ID", ok)
		err = common.ErrorUserNotLogin
		return
	}
	return
}
