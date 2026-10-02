package controller

import (
	"newaoe/Src/util"
	"newaoe/Src/vip/service"

	"github.com/gin-gonic/gin"
)

func CancelSubmitBlock(ctx *gin.Context) {
	err := service.SetBlockAllSubmit(false)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "取消阻塞提交成功!", nil)
}
