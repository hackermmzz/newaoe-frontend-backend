package controller

import (
	"newaoe/Src/codeRun/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

type CodeRunInfoFromPostBost struct {
	Indices int `json:"indices"`
	Class   int `json:"class"`
	RunType int `json:"runType"`
}

func CodeRun(ctx *gin.Context) {
	//
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//解析数据
	var info CodeRunInfoFromPostBost
	if !util.JsonCtx(ctx, &info) {
		util.DebugError("CodeRun解析数据失败!")
		util.ResponseNAK_MSG(ctx, "服务器失败!", nil)
		return
	}
	//运行代码
	err := service.CodeRun(info.Indices, id, info.Class, info.RunType)
	if err != nil {
		util.DebugError("CodeRun", err)
		util.ResponseNAK_MSG(ctx, "服务器失败!", nil)
		return
	}
	//
	util.DebugSuccess("代码运行成功!")
	util.ResponseACK_MSG(ctx, "运行成功!", nil)
}
