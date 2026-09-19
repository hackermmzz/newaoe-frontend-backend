package controller

import (
	"newaoe/Src/codeRun/service"
	"newaoe/Src/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func StudentHistoryGet(ctx *gin.Context) {
	//获取用户数据
	data := util.GetCtxTookenInfo(ctx)
	if data == nil {
		util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
		return
	}
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
	//获取查寻的id
	student_id := ctx.Query("student_id")
	id := data["id"].(string)
	if student_id != "" {
		id = student_id
	}
	//获取历史记录
	submitRecords, err := service.GetHistoryRangeById(id, beg, end)
	if err != nil {
		util.DebugError("StudentHistoryGet", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//
	util.ResponseACK_MSG(ctx, "历史记录获取成功", submitRecords)
}
