package controller

import (
	"newaoe/Src/user/dao"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func StudentInfoGet(ctx *gin.Context) {
	dt := util.GetCtxTookenInfo(ctx)
	if dt == nil {
		util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
		return
	}
	//获取id
	id := dt["id"].(string)
	//返回数据
	usr := dao.UserGet(id)
	util.ResponseACK_MSG(ctx, "获取成功!", usr)
}
