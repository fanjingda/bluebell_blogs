package mysql

import (
	"bluebell_blogs/models/common"
	"bluebell_blogs/models/entity"
	"bluebell_blogs/pkg/crypto_md5"
	"database/sql"
	"fmt"
)

func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return common.ErrorUserExist
	}
	return
}

func InsertUser(user *entity.User) (err error) {
	user.Password = crypto_md5.EncyptPassword(user.Password)
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	if err != nil {
		fmt.Println("插入错误", err)
	}
	return
}

func Login(user *entity.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id,username,password from user where username=?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return common.ErrorUserNotExist
	}
	if err != nil {
		return err
	}
	password := crypto_md5.EncyptPassword(oPassword)
	if password != user.Password {
		return common.ErrorInvalidPassword
	}
	return
}

func GetUserByID(uid int64) (user *entity.User, err error) {
	user = new(entity.User)
	sqlStr := `select user_id username from user where user_id=?`
	err = db.Get(user, sqlStr, uid)
	return
}
