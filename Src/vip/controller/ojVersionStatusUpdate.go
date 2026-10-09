package controller

import (
	"newaoe/Src/util"
	"newaoe/Src/vip/service"

	"github.com/gin-gonic/gin"
)

func OJVersionStatusUpdate(ctx *gin.Context) {
	latestVersion := ctx.Query("latestVersion")
	//
	err := service.OJVersionUpdate(latestVersion)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "版本推送成功!", nil)
}
