package controller

import (
	"newaoe/Src/user/dao"
	"newaoe/Src/user/model"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func EnhanceVip(ctx *gin.Context) {
	//获取user
	id := ctx.Query("id")
	if id == "" {
		util.ResponseNAK_MSG(ctx, "id不能为空", "")
		return
	}
	//更新vip等级
	ok, err := dao.UserUpdate(nil, id, model.Student{
		Vip: model.VIP_SUPER,
	})
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	if !ok {
		util.ResponseNAK_MSG(ctx, "更新失败", nil)
		return
	}
	util.ResponseACK_MSG(ctx, "更新成功", "")
}
