package controller

import (
	"newaoe/Src/map/dao"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func DownloadLatestAOEPackage(ctx *gin.Context) {
	//直接查数据库
	value, err := dao.MapGet(nil, "newaoe-latest-version")
	if err != nil {
		util.DebugError("MapGet:", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "获取下载链接成功!", map[string]interface{}{"url": value})
}
