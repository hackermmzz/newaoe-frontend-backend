package controller

import (
	"newaoe/Src/user/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

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
	//获取用户信息
	info, err := service.StudentInfoGetByEmailOrID(dt.Id)
	if err != nil {
		util.DebugError("重置密码失败:", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//发送验证码
	err = service.PasswordResetCodeSend(info.Id, info.Email)
	if err != nil {
		util.DebugError("密码重置验证码发送失败!")
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//
	util.ResponseACK_MSG(ctx, "验证码发送成功", map[string]interface{}{"email": dt.Email})
}
