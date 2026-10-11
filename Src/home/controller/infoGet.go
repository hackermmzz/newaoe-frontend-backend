package controller

import (
	"errors"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func StudentInfoGet(ctx *gin.Context) {
	dt := util.GetCtxTookenInfo(ctx)
	if dt == nil {
		util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
		return
	}
	//获取id
	id := dt["id"].(string)
	//返回数据
	usr, err := dao.UserGet(nil, id)
	if err != nil || usr == nil {
		if usr == nil {
			err = errors.New("用户不存在")
		}
		util.DebugError("UserGet:", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	usr.Password = "empty"
	util.ResponseACK_MSG(ctx, "获取成功!", usr)
}
