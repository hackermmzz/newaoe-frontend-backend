package controller

import (
	"newaoe/Src/user/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

type touristRegistCodeSendDataInfo struct {
	Email string `json:"email"`
}

func TouristRegistCodeSend(ctx *gin.Context) {
	var dt touristRegistCodeSendDataInfo
	if !util.JsonCtx(ctx, &dt) {
		util.ResponseNAK_MSG(ctx, "数据报文错误", "")
		return
	}
	//发送验证码(这里把邮箱当作redis的id就行了)
	err := service.SendRegistCode(dt.Email, dt.Email)
	if err != nil {
		util.DebugError("发送验证码失败", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//
	util.ResponseACK_MSG(ctx, "验证码发送成功", map[string]interface{}{"email": dt.Email})
}
