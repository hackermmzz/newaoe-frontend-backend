package controller

import (
	"errors"
	"fmt"
	"newaoe/Src/user/dao"
	"newaoe/Src/user/model"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func EnhanceVip(ctx *gin.Context) {
	//获取user
	idOrEmail := ctx.Query("id")
	if idOrEmail == "" {
		util.ResponseNAK_MSG(ctx, "id不能为空", "")
		return
	}
	//获取用户的id
	student, err := dao.UserGetByIdOrEmail(nil, idOrEmail)
	if err != nil || student == nil {
		if err == nil {
			err = errors.New(fmt.Sprintf("用户不存在! idOrEmail: %s", idOrEmail))
		}
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//更新vip等级
	ok, err := dao.UserUpdate(nil, student.ID, model.Student{
		Vip: model.VIP_SUPER,
	})
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	if !ok {
		util.ResponseNAK_MSG(ctx, "权限更新失败", nil)
		return
	}
	util.ResponseACK_MSG(ctx, "权限更新成功", "")
}
