package controller

import (
	"newaoe/Src/home/dao"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func FetchAnnouncement(ctx *gin.Context) {
	//直接获取最新的版本
	data, err := dao.AnnouncementGetLatest(nil)
	if err != nil {
		util.DebugError("AnnouncementGetLatest:", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	if data == nil {
		util.ResponseACK_MSG(ctx, "暂无数据!", nil)
		return
	}
	util.ResponseACK_MSG(ctx, "获取成功!", data)
}
