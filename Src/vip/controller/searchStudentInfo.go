package controller

import (
	"newaoe/Src/user/dao"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func SearchStudentInfo(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		util.ResponseNAK_MSG(ctx, "格式错误!", nil)
		return
	}
	//获取学生信息
	data := dao.UserGet(nil, id)
	util.ResponseACK_MSG(ctx, "搜索成功!", data)
}
