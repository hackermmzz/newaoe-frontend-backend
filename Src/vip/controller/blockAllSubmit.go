package controller

import (
	"newaoe/Src/util"
	"newaoe/Src/vip/service"

	"github.com/gin-gonic/gin"
)

func BlockAllSubmit(ctx *gin.Context) {
	err := service.SetBlockAllSubmit(true)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "阻塞所有提交成功!", nil)
}
