package controller

import (
	"newaoe/Src/util"
	"newaoe/Src/vip/service"

	"github.com/gin-gonic/gin"
)

func OJVersionStatusUpdate(ctx *gin.Context) {
	latestVersion := ctx.Query("latestVersion")
	hard_reboot := ctx.Query("hard_reboot")
	//
	err := service.OJVersionUpdate(latestVersion, hard_reboot == "true")
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "版本推送成功!", nil)
}
