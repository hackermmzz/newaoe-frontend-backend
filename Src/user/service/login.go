package service

import (
	"errors"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
)

// id和加密前的密码
func UserCanLogin(id string, password string) error {
	correct_password := dao.UserGetPassword(nil, id)
	if !util.CheckPasswordSame(correct_password, password) {
		return errors.New("账号或者密码错误!")
	}
	return nil
}
