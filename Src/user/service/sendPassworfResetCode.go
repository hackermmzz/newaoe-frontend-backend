package service

import (
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/email/service"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
)

func PasswordResetCodeSend(id string, email string) error {
	//判断是否可以发送验证码
	ttl, err := dao.PasswordForgetVerifyCodeCanSendTTL(id)
	if err != nil {
		return err
	}
	if ttl > 0 {
		return util.NewError(fmt.Sprintf("请等待%v秒后重试!", ttl))
	}
	//发送验证码
	code := util.GenerateVerifyCode(config.Conf.PasswordForgetVerifyCode.Length)
	emailmsg := service.EmailMsg{
		Email:   email,
		Subject: "重置密码验证码",
		Text:    "你的验证码为: " + code,
		Type:    service.EmailMsgType_TEXT,
	}
	service.SendEmail(emailmsg)
	//插入Redis
	if err := dao.PasswordForgetVerifyCodeAdd(id, code); err != nil {
		return err
	}
	//
	return nil
}
