package controller

import (
	"newaoe/Src/user/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	//数据结构体
	type DataInfo struct {
		Id       string `json:"id"` //前端传过来的可能是id也可能是邮箱，这里统一叫id
		Password string `json:"password"`
	}
	//
	var dt DataInfo
	if !util.JsonCtx(c, &dt) {
		util.ResponseNAK_MSG(c, "数据报文错误", "")
		return
	}
	//根据传入的账号/邮箱获取学生的id
	id, err := service.StudentIDGetByEmailOrID(dt.Id)
	if err != nil {
		util.DebugError("UserLogin:", err)
		util.ResponseNAK_MSG(c, err.Error(), nil)
		return
	}
	dt.Id = id
	//登陆
	if err := service.UserCanLogin(dt.Id, dt.Password); err != nil {
		util.DebugError("UserLogin:", err)
		util.ResponseNAK_MSG(c, err.Error(), nil)
		return
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
