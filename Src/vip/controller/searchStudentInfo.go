package controller

import (
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func SearchStudentInfo(ctx *gin.Context) {
	//获取id_like
	id := ctx.Query("id")
	if id == "" {
		util.ResponseNAK_MSG(ctx, "格式错误!", nil)
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
	//获取学生信息
	data, err := dao.UserGetLikeId(nil, id, beg, end+1)
	if err != nil {
		util.DebugError("UserGetLikeId:", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "搜索成功!", data)
}
