package controller

import (
	"newaoe/Src/home/dao"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func FetchAnnouncement(ctx *gin.Context) {
	//直接获取最新的版本
	data := dao.AnnouncementGetLatest(nil)
	if data == nil {
		util.ResponseACK_MSG(ctx, "暂无数据!", nil)
		return
	}
	util.ResponseACK_MSG(ctx, "获取成功!", data)
}
