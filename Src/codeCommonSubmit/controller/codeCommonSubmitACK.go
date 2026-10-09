package controller

import (
	"newaoe/Src/codeCommonSubmit/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

type codeCommonPostBodyInfo struct {
	Key string `json:"key"`
}

func CodeCommonSubmitACK(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//获取postBody数据
	var info codeCommonPostBodyInfo
	if !util.JsonCtx(ctx, &info) {
		util.DebugError("CodeCommonSubmitACK解析PostBody失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//获取记录
	indices, err := service.CodeCommonSubmitACK(id, info.Key)
	if err != nil {
		util.DebugError("CodeCommonSubmitACK", err)
		util.ResponseNAK_MSG(ctx, "服务器异常", nil)
		return
	}
	//
	util.ResponseACK_MSG(ctx, "代码上传成功", map[string]interface{}{"indices": indices})
	util.DebugSuccess(id, "代码上传成功!")
}
