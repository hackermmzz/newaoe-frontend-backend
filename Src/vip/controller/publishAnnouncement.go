package controller

import (
	"newaoe/Src/util"
	"newaoe/Src/vip/service"

	"github.com/gin-gonic/gin"
)

type announcementPostDataInfo struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func PublishAnnouncement(ctx *gin.Context) {
	var info announcementPostDataInfo
	if !util.JsonCtx(ctx, &info) {
		util.ResponseNAK_MSG(ctx, "格式错误!", nil)
		return
	}
	//默认都是enabled
	err := service.PublishAnnouncement(info.Title, info.Content, true)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "发布成功!", nil)
}
