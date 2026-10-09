package controller

import (
	"fmt"
	"newaoe/Src/util"
	"newaoe/Src/vip/service"

	"github.com/gin-gonic/gin"
)

func ReRunAnomalRecord(ctx *gin.Context) {
	totalCnt, err := service.ReRunAnomalRecord()
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//
	msg := fmt.Sprintf("重新运行成功!记录有%d条", totalCnt)
	util.DebugSuccess(msg)
	util.ResponseACK_MSG(ctx, msg, nil)
}
