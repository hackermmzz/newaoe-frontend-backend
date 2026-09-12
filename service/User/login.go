package User

import (
	data "newaoe/service/Data"
	"newaoe/util"

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
	//获取用户密码
	correct_password := data.UserGetPassword(dt.Id)
	if !util.CheckPasswordSame(correct_password, dt.Password) {
		util.ResponseNAK_MSG(c, "账号或者密码错误!", "")
		return
	}
	//创建cookie
	token := UserSetCookie(c, dt.Id)
	//写入数据库
	if token == "" {
		util.ResponseNAK_MSG(c, "服务器异常", "")
		return
	}
	//
	util.ResponseACK_MSG(c, "登陆成功", "")
}
