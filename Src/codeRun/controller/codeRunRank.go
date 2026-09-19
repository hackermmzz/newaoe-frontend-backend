package controller

import (
	"newaoe/Src/codeRun/service"
	"newaoe/Src/util"

	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func FetchRank(ctx *gin.Context) {
	//解析参数
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
	//获取beg到end的排行数据
	data, err := service.RankFetch(beg, end+1)
	if err != nil {
		util.DebugError("FetchRank", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//返回结果
	util.ResponseACK_MSG(ctx, "获取排行成功", data)
}
