package controller

import (
	"fmt"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/email/service"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

// 重置密码
func UserResetPassword(ctx *gin.Context) {
	type DataInfo struct {
		Id         string `json:"id"`
		Password   string `json:"password"`
		Email      string `json:"email"`
		Verifycode string `json:"verifycode"`
	}
	var dt DataInfo
	if !util.JsonCtx(ctx, &dt) {
		util.ResponseNAK_MSG(ctx, "数据报文错误", "")
		return
	}
	//处理邮箱
	dt.Email = userEmailProcess(dt.Id, dt.Email)
	//检查用户名是否存在
	if !dao.UserExist(dt.Id) {
		util.ResponseNAK_MSG(ctx, "用户不存在", "")
		return
	}
	//检查用户名是否与用户邮箱匹配
	if dao.UserGetEmail(dt.Id) != dt.Email {
		util.ResponseNAK_MSG(ctx, "用户名与邮箱不匹配", "")
		return
	}
	//
	//检查验证码是否正确
	if !dao.PasswordForgetVerifyCodeExist(dt.Id, dt.Verifycode) {
		util.ResponseNAK_MSG(ctx, "验证码错误", "")
		return
	}
	//检查密码是否符合格式
	if !PasswordLegal(dt.Password) {
		util.ResponseNAK_MSG(ctx, "密码格式错误!", "")
		return
	}
	//检查验证码是否正确
	if !dao.PasswordForgetVerifyCodeExist(dt.Id, dt.Verifycode) {
		util.ResponseNAK_MSG(ctx, "验证码错误", "")
		return
	}
	//重置密码
	session := database.NewSession()
	defer session.Rollback()
	err := session.Begin()
	if err != nil {
		util.ResponseNAK_MSG(ctx, "服务器异常", "")
		return
	}

	password, err := util.EncodePassword(dt.Password)
	if err != nil {
		util.ResponseNAK_MSG(ctx, "密码格式不合规则!", "")
		return
	}
	if !dao.UserResetPassword(session, dt.Id, password) {
		util.ResponseNAK_MSG(ctx, "重置密码失败", "")
		return
	}

	if err = session.Commit(); err != nil {
		util.ResponseNAK_MSG(ctx, "服务器异常!", "")
		return
	}
	//注册成功
	util.ResponseACK_MSG(ctx, "密码更改成功", "")
}

// 重置密码验证码
func UserPasswordForgetCodeSend(ctx *gin.Context) {
	type DataInfo struct {
		Id    string `json:"id"`
		Email string `json:"email"`
	}
	//
	var dt DataInfo
	if !util.JsonCtx(ctx, &dt) {
		util.ResponseNAK_MSG(ctx, "数据报文错误", "")
		return
	}
	//处理邮箱
	dt.Email = userEmailProcess(dt.Id, dt.Email)
	//判断是否可以发送验证码
	ttl := dao.PasswordForgetVerifyCodeCanSendTTL(dt.Id)
	if ttl > 0 {
		util.ResponseNAK_MSG(ctx, fmt.Sprintf("请等待%v秒后重试!", ttl), "")
		return
	}
	//发送验证码
	code := util.GenerateVerifyCode(config.Conf.PasswordForgetVerifyCode.Length)
	email := service.EmailMsg{
		Email:   dt.Email,
		Subject: "重置密码验证码",
		Text:    "你的验证码为: " + code,
		Type:    service.EmailMsgType_TEXT,
	}
	service.SendEmail(email)
	//插入Redis
	success := dao.PasswordForgetVerifyCodeAdd(dt.Id, code)
	//
	if !success {
		util.ResponseNAK_MSG(ctx, "验证码发送失败", "")
		return
	}
	util.ResponseACK_MSG(ctx, "验证码发送成功", "")
}
