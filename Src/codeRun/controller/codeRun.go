package controller

import (
	"context"
	"encoding/json"
	codeAssessmentDao "newaoe/Src/codeAssessmentSubmit/dao"
	codeCommonDao "newaoe/Src/codeCommonSubmit/dao"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type CodeRunInfoFromPostBost struct {
	Indices int `json:"indices"`
	Class   int `json:"class"`
}

func CodeRun(ctx *gin.Context) {
	//
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//
	session := database.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		util.DebugError(id, "代码运行失败", err.Error())
		util.ResponseNAK_MSG(ctx, "代码运行失败", "")
		return
	}
	defer session.Rollback()
	//解析数据
	var info CodeRunInfoFromPostBost
	if !util.JsonCtx(ctx, &info) {
		util.DebugError("CodeRun解析数据失败!")
		util.ResponseNAK_MSG(ctx, "服务器失败!", nil)
		return
	}
	//获取数据库数据
	codeRunInfo := model.CodeRunInfo{
		ID:         id,
		SubmitTime: util.UTC_Time(),
		Class:      info.Class,
		Status:     model.ProcessDataMessageByStatus(model.NewCodeRunStatusInfo()).Marshal(),
	}
	switch info.Class {
	case model.Code_CommonSubmit:
		//查表拿到数据
		data := codeCommonDao.CodeCommonGetByIndices(info.Indices)
		if data == nil || data.ID != id {
			util.DebugError("有人作假", id)
			util.ResponseNAK_MSG(ctx, "服务器异常", nil)
			return
		}
		codeRunInfo.Header = data.Header
		codeRunInfo.Source = data.Source
		codeRunInfo.Description = data.Description
	case model.Code_AssessmentSubmit:
		data := codeAssessmentDao.CodeAssessmentGetByIndices(info.Indices)
		if data == nil || data.ID != id {
			util.DebugError("有人作假", id)
			util.ResponseNAK_MSG(ctx, "服务器异常", nil)
			return
		}
		codeRunInfo.Header = data.Header
		codeRunInfo.Source = data.Source
		codeRunInfo.Description = ""
	case model.Code_ReRunSubmit:
		data := dao.CodeRunGetByIndices(info.Indices)
		if data == nil || data.ID != id {
			util.DebugError("有人作假", id)
			util.ResponseNAK_MSG(ctx, "服务器异常", nil)
			return
		}
		codeRunInfo.Header = data.Header
		codeRunInfo.Source = data.Source
		codeRunInfo.Description = data.Description
	default:
		util.DebugError("用户传入了一个垃圾类型!")
		util.ResponseNAK_MSG(ctx, "服务器异常", nil)
		return
	}
	//向数据库插入一条运行记录
	codeRunInfo.Indices = dao.CodeRunAdd(session, codeRunInfo)
	if codeRunInfo.Indices == 0 {
		util.DebugError("RunUserCode插入CodeRunAdd记录失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//插入一条记录（以防止崩溃可以恢复）
	if !dao.CodeRunningInsert(session, model.CodeRunningInfo{Indices: codeRunInfo.Indices, SubmitTime: codeRunInfo.SubmitTime}) {
		util.DebugError("RunUserCode插入CodeRunningInsert记录失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//插入队列，准备给oj消费
	if !addCodeFile(codeRunInfo) {
		//
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	if err := session.Commit(); err != nil {
		util.DebugError("CodeRun提交session记录失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	util.ResponseACK_MSG(ctx, "运行成功!", nil)
}

// 根据索引将对应的代码文件加入待运行队列
func RunUserCode(session *xorm.Session, info model.CodeRunInfo, needInsert bool) bool {
	if needInsert {

	} else {
		//清空原先状态
		status := model.NewCodeRunStatusInfo()
		if !dao.CodeRunUpdateStatusAsync(session, info.Indices, status.Marshal()) {
			util.DebugError("RunUserCode update status fail!")
			return false
		}
	}

	return true
}

// 代码文件加入待运行队列
func addCodeFile(task model.CodeRunInfo) bool {
	data, err := json.Marshal(task)
	if err != nil {
		util.Debug("AddCodeFile json.Marshal err:", err)
		return false
	}
	return redis.RedisListPush(context.Background(), config.Conf.Code.CodeWaitForRunQueueTopic, string(data))
}
