package controller

import (
	"newaoe/Src/codeRun/service"
	"newaoe/Src/util"
	"time"

	"github.com/gin-gonic/gin"
)

func ReRunAnomalRecord(ctx *gin.Context) {
	expireDuration := time.Duration(-1) * time.Second
	for {
		cnt, err := service.ReRunLongTimeWaitRecord(expireDuration)
		if err != nil {
			util.DebugError("ReRunAnomalRecord", err)
			util.ResponseNAK_MSG(ctx, err.Error(), nil)
			return
		}
		if cnt == 0 {
			break
		}
	}
	//

	util.ResponseACK_MSG(ctx, "重新运行成功!", nil)
}
