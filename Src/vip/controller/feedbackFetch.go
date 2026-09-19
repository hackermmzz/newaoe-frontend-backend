package controller

import (
	"newaoe/Src/util"
	"newaoe/Src/vip/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func FetchFeedback(ctx *gin.Context) {
	//获取范围
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
	//获取反馈记录
	feebackRecord, err := service.FetchFeedbackRecord(beg, end+1)
	if err != nil {
		util.DebugError("FetchFeedback", err)
		util.ResponseNAK_MSG(ctx, "服务器异常!", "")
		return
	}
	//
	util.ResponseACK_MSG(ctx, "反馈记录获取成功", feebackRecord)
}
