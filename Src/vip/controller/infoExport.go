package controller

import (
	"newaoe/Src/util"
	"newaoe/Src/vip/service"

	"github.com/gin-gonic/gin"
)

func StudentInfoExport(ctx *gin.Context) {
	url, err := service.ExcelInfoExport()
	if err != nil {
		util.DebugError("StudentInfoExport:", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//返回路径
	util.ResponseACK_MSG(ctx, "导出成功!", url)
}
