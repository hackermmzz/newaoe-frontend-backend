package controller

import (
	"newaoe/Src/user/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	//数据结构体
	type DataInfo struct {
		Id       string `json:"id"`
		Password string `json:"password"`
	}
	//
	var dt DataInfo
	if !util.JsonCtx(c, &dt) {
		util.ResponseNAK_MSG(c, "数据报文错误", "")
		return
	}
	//登陆
	if err := service.UserCanLogin(dt.Id, dt.Password); err != nil {
		util.Debug("UserLogin:", err)
		util.ResponseNAK_MSG(c, err.Error(), nil)
	}
	//设置cookie
	cookie, err := service.UserGenCookie(dt.Id, c.ClientIP())
	if err != nil {
		util.DebugError("生成cookie失败:", err)
		util.ResponseNAK_MSG(c, "服务器异常!", "")
		return
	}
	c.Header("set-cookie", cookie)
	//
	util.ResponseACK_MSG(c, "登陆成功", "")
}
