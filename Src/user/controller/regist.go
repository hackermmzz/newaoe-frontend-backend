package controller

import (
	"newaoe/Src/user/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func UserRegist(c *gin.Context) {
	//数据结构体
	type DataInfo struct {
		Id         string `json:"id"`
		Password   string `json:"password"`
		Email      string `json:"email"`
		Verifycode string `json:"verifycode"`
	}
	//
	var dt DataInfo
	if !util.JsonCtx(c, &dt) {
		util.ResponseNAK_MSG(c, "数据报文错误", "")
		return
	}
	//处理邮箱
	dt.Email = userEmailProcess(dt.Id, dt.Email)
	//注册
	err := service.UserRegist(dt.Id, dt.Password, dt.Email, dt.Verifycode)
	if err != nil {
		util.DebugError("UserRegist:", err)
		util.ResponseNAK_MSG(c, err.Error(), nil)
		return
	}
	//注册成功
	util.ResponseACK_MSG(c, "注册成功", "")
}
