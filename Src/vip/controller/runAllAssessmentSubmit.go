package controller

import (
	"fmt"
	"newaoe/Src/util"
	"newaoe/Src/vip/service"

	"github.com/gin-gonic/gin"
)

func RunAllAssessmentSubmit(ctx *gin.Context) {
	cnt, err := service.RunAllAssessmentSubmit()
	if err != nil {
		util.DebugError("RunAllAssessmentSubmit", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, fmt.Sprintf("%d条记录全部运行成功!", cnt), nil)
}
