package controller

import (
	"newaoe/Src/codeCommonSubmit/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func CodeCommonSubmit(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//获取下载链接和数据key
	urls, key, err := service.CodeCommonSubmit(id)
	if err != nil {
		util.DebugError("CodeCommonSubmit:", err)
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//回复
	util.ResponseACK_MSG(ctx, "获取链接成功!",
		map[string]interface{}{
			"urls": urls,
			"key":  key,
		})
}
