package Home

import (
	data "newaoe/service/Data"
	"newaoe/util"

	"github.com/gin-gonic/gin"
)

func StudentInfoGet(ctx *gin.Context) {
	dt := util.GetCtxTookenInfo(ctx)
	if dt == nil {
		util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
		return
	}
	//获取数据
	email := dt["email_bind"]
	regist_date := dt["regist_date"]
	id := dt["id"]
	//返回数据
	util.ResponseACK_MSG(ctx, "获取成功!", map[string]string{
		"id":          id,
		"email":       email,
		"avatar":      data.UserGetAvatar(id), //后端转发
		"regist_date": regist_date,
	})
}
