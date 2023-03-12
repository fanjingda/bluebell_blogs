package logic

import (
	"bluebell_blogs/dao/mysql"
	"bluebell_blogs/models/dto"
	"bluebell_blogs/models/entity"
	"bluebell_blogs/pkg/jwt"
	"bluebell_blogs/pkg/snowflake"
	"fmt"
)

func SignUp(p *dto.ParamSignUp) (err error) {
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		fmt.Println(" mysql.CheckUserExist;用户错误", err)
		return err
	}
	userID := snowflake.GenID()

	user := &entity.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.InsertUser(user)
}

func Login(p *dto.ParamLogin) (token string, err error) {
	user := &entity.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		fmt.Println("Login 连接失败！", err)
		return "", err
	}

	return jwt.GenToken(user.UserID, user.Username)
}
