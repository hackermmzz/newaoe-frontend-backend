package controller

import (
	"newaoe/Src/user/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

type passwordResetDataInfo struct {
	Id         string `json:"id"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Verifycode string `json:"verifycode"`
}

// 重置密码
func UserResetPassword(ctx *gin.Context) {

	var dt passwordResetDataInfo
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
	//
	err = service.PasswordReset(info.Id, info.Email, dt.Password, dt.Verifycode)
	if err != nil {
		util.DebugError("重置密码失败:", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//更改成功
	util.ResponseACK_MSG(ctx, "密码更改成功", "")
}
