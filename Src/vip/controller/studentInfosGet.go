package controller

import (
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
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
	if err := util.ValidateInclusiveRange(beg, end); err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//获取数据
	dt, err := dao.UserGetByRangeOrderByID(nil, beg, end+1)
	if err != nil {
		util.DebugError("UserGetByRangeOrderByID:", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//隐藏密码
	for i := range dt {
		dt[i].Password = "empty"
	}
	//
	util.ResponseACK_MSG(ctx, "获取成功!", dt)
}
