package service

import (
	"errors"
	database "newaoe/Src/databse"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
)

func PasswordReset(id string, targetEmail string, newPasword string, verifyCode string) error {
	//
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return errors.New("事务开启失败!")
	}
	defer session.Rollback()
	//检查用户名是否存在
	exist, err := dao.UserExist(session, id)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("用户不存在")
	}
	//检查用户名是否与用户邮箱匹配
	userEmail, err := dao.UserGetEmail(session, id)
	if err != nil {
		return err
	}
	if userEmail != targetEmail {
		return errors.New("用户名与邮箱不匹配")
	}
	//
	//检查验证码是否正确
	ok, err := dao.PasswordForgetVerifyCodeExist(id, verifyCode)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("验证码错误")
	}
	//移除验证码
	if _, err := dao.PasswordForgetVerifyCodeRemove(id); err != nil {
		return err
	}
	//检查密码是否符合格式
	if !passwordLegal(newPasword) {
		return errors.New("密码格式错误!")
	}
	//重置密码
	newEncodePassword, err := util.EncodePassword(newPasword)
	if err != nil {
		return errors.New("密码格式不合规则!")
	}
	updated, err := dao.UserResetPassword(session, id, newEncodePassword)
	if err != nil {
		return err
	}
	if !updated {
		return errors.New("重置密码失败")
	}
	//提交业务
	if err = session.Commit(); err != nil {
		return errors.New("服务器异常!")
	}
	return nil
}
