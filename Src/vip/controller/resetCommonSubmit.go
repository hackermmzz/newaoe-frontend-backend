package controller

import (
	"newaoe/Src/util"
	"newaoe/Src/vip/service"

	"github.com/gin-gonic/gin"
)

type resetCommonSubmitPostBody struct {
	StudentIDS []string `json:"studentIDs"`
}

func ResetCommonSubmit(ctx *gin.Context) {
	//获取学生id
	var data resetCommonSubmitPostBody
	if !util.JsonCtx(ctx, &data) {
		util.ResponseNAK_MSG(ctx, "数据格式异常!", nil)
		return
	}
	//
	err := service.ResetCommonSubmit(data.StudentIDS)
	if err != nil {
		util.DebugError("ResetCommonSubmit", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	util.ResponseACK_MSG(ctx, "重置成功!", nil)
}
