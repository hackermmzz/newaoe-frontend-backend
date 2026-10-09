package controller

import (
	"newaoe/Src/codeRun/model"
	"newaoe/Src/codeRun/service"
	"newaoe/Src/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func FetchSubmitRecord(ctx *gin.Context) {
	range_ := ctx.Query("range")
	parts := strings.Split(range_, ":")
	if len(parts) != 2 {
		util.ResponseNAK_MSG(ctx, "参数传递错误!", "")
		return
	}
	beg, err0 := strconv.Atoi(parts[0])
	end, err1 := strconv.Atoi(parts[1])
	if err0 != nil || err1 != nil {
		util.ResponseNAK_MSG(ctx, "参数传递错误!", "")
		return
	}
	//
	data, err := service.GetHistoryRangeById("", beg, end, model.Code_CommonSubmit, model.Code_ReRunSubmit)
	if err != nil {
		util.DebugError("FetchSubmitRecord", err)
		util.ResponseNAK_MSG(ctx, "获取失败!", err.Error())
		return
	}
	util.ResponseACK_MSG(ctx, "获取成功!", data)
}
