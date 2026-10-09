package controller

import (
	"newaoe/Src/home/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func AvatarUpdate(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id, _ := userInfo["id"].(string)
	//生成上传链接
	url, err := service.AvatarUpdate(id)
	if err != nil {
		util.DebugError("AvatarUpdate", err)
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	util.ResponseACK_MSG(ctx, "获取上传链接成功!", map[string]interface{}{"url": url})
}
