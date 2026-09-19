package service

import (
	"errors"
	database "newaoe/Src/databse"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
)

func UserRegist(id string, password string, email string, verifyCode string) error {
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return errors.New("事务开启失败!")
	}
	defer session.Rollback()
	//检查密码是否符合格式
	if !passwordLegal(password) {
		return errors.New("密码格式错误!")
	}
	//检查验证码是否正确
	if !dao.RegistVerifyCodeExist(id, verifyCode) {
		return errors.New("验证码错误")
	}
	//查询用是否允许注册
	if !userIdLegal(id) {
		return errors.New("该用户不允许注册账号!")
	}
	//查询是否已经注册
	if dao.UserExist(session, id) {
		return errors.New("用户已经注册")
	}
	//查询邮箱是否已经使用过
	if dao.EmailExist(session, email) {
		return errors.New("邮箱已经使用")
	}
	//添加用户
	password, err := util.EncodePassword(password)
	if err != nil {
		return errors.New("密码格式不合规则!")
	}
	if !dao.UserAdd(session, id, password, email) {
		return errors.New("注册失败")
	}

	if err = session.Commit(); err != nil {
		return errors.New("服务器异常!")
	}
	return nil
}

// 判断密码是否合法
func passwordLegal(password string) bool {
	return true
}

// 判断用户是否允许注册(不允许未授权用户注册)
func userIdLegal(id string) bool {
	return true
}
