package controller

import (
	"newaoe/Src/user/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

type registCodeSendDataInfo struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

// 常规注册 发送注册码
func UserRegistCodeSend(ctx *gin.Context) {
	//
	var dt registCodeSendDataInfo
	if !util.JsonCtx(ctx, &dt) {
		util.ResponseNAK_MSG(ctx, "数据报文错误", "")
		return
	}
	//处理邮箱
	dt.Email = userEmailProcess(dt.Id, dt.Email)
	//发送验证码
	err := service.SendRegistCode(dt.Id, dt.Email)
	if err != nil {
		util.DebugError("发送验证码失败", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//
	util.ResponseACK_MSG(ctx, "验证码发送成功", map[string]interface{}{"email": dt.Email})
}

// 对用户绑定的邮箱进行处理
func userEmailProcess(id string, email string) string {
	//强制邮箱只能是Id@njust.edu.cn
	return id + "@njust.edu.cn"
}
