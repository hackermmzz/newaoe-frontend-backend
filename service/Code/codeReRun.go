package Code

import (
	"newaoe/dao"
	"newaoe/util"

	"github.com/gin-gonic/gin"
)

// 代码重新运行
func CodeReRun(ctx *gin.Context) {
	//获取用户信息
	userInfo := util.GetCtxTookenInfo(ctx)
	if userInfo == nil {
		util.ResponseNAK_MSG(ctx, "运行失败", "")
		return
	}
	//
	id := userInfo["id"]
	//获取对应的文件
	type MSG struct {
		Source      string `json:"source"`
		Header      string `json:"header"`
		Description string `json:"description"`
	}
	var msg MSG
	if !util.JsonCtx(ctx, &msg) {
		util.ResponseNAK_MSG(ctx, "运行失败", "")
		return
	}
	//将代码加入代运行就绪队列
	session := dao.DB.NewSession()
	defer session.Rollback()
	err := session.Begin()
	if err != nil {
		util.ResponseNAK_MSG(ctx, "服务器异常!", "")
		return
	}
	//运行代码
	runInfo := dao.CodeRunInfo{
		ID:          id.(string),
		SubmitTime:  util.UTC_Time(),
		Header:      msg.Header,
		Source:      msg.Source,
		Class:       dao.Code_Common,
		Description: msg.Description,
		Status:      dao.ProcessDataMessageByStatus(dao.NewCodeRunStatusInfo()).Marshal(),
		Version:     util.UTC_Time().UnixMilli(),
	}
	if !RunUserCode(session, runInfo, true) {
		util.ResponseNAK_MSG(ctx, "运行失败!", "")
		return
	}
	//提交
	if err = session.Commit(); err != nil {
		util.ResponseNAK_MSG(ctx, "服务器异常!", "")
		return
	}
	//
	util.ResponseACK_MSG(ctx, "正在运行", "")
}
