package service

import (
	"errors"
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/email/service"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
)

func SendRegistCode(id string, email string) error {
	//判断是否可以发送验证码
	ttl := dao.RegistVerifyCodeCanSendTTL(id)
	if ttl > 0 {
		return errors.New(fmt.Sprintf("请等待%v秒后重试!", ttl))
	}
	//发送验证码
	code := util.GenerateVerifyCode(config.Conf.RegistVerifyCode.Length)
	emailMsg := service.EmailMsg{
		Email:   email,
		Subject: "注册验证码",
		Text:    "你的验证码为: " + code,
		Type:    service.EmailMsgType_TEXT,
	}
	service.SendEmail(emailMsg)
	//插入Redis
	dao.RegistVerifyCodeAdd(id, code)
	//
	return nil
}
