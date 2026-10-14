package controller

import (
	"fmt"
	"newaoe/Src/forum/service"
	"newaoe/Src/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ForumAnnouncementRemove(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	vip := int(userInfo["vip"].(float64))
	//获取要移除的索引
	indicesStr := ctx.Query("indices")
	indices, err := strconv.ParseInt(indicesStr, 10, 64)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	err = service.ForumAnnouncementRemove(indices, id, vip)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, fmt.Sprintf("帖子%d删除成功!", indices), nil)
}
