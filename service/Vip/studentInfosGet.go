package vip

import (
	data "newaoe/service/Data"
	"newaoe/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func StudentInfosGet(ctx *gin.Context) {
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
	//获取数据
	dt := data.UserGetByRangeOrderByRegistData(beg, end+1)
	//
	util.ResponseACK_MSG(ctx, "获取成功!", dt)
}
