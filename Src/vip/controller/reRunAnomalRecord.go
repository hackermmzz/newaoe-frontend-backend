package controller

import (
	"errors"
	"fmt"
	"newaoe/Src/codeRun/service"
	"newaoe/Src/util"
	"time"

	"github.com/gin-gonic/gin"
)

func ReRunAnomalRecord(ctx *gin.Context) {
	expireDuration := time.Duration(-24) * time.Hour
	totalCnt := 0
	errs := make([]error, 0)
	for {
		cnt, err := service.ReRunLongTimeWaitRecord(expireDuration)
		if err != nil {
			errs = append(errs, err)
		}
		totalCnt += cnt
		if cnt == 0 {
			break
		}
	}
	//
	if len(errs) != 0 {
		err := errors.Join(errs...).Error()
		util.DebugError("ReRunAnomalRecord", err)
		util.ResponseNAK_MSG(ctx, err, nil)
		return
	}
	//
	msg := fmt.Sprintf("重新运行成功!记录有%d条", totalCnt)
	util.DebugSuccess(msg)
	util.ResponseACK_MSG(ctx, msg, nil)
}
