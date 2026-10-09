package controller

import (
	"newaoe/Src/codeAssessmentSubmit/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

type assessmentACKBody struct {
	Key            string `json:"key"`
	Teacher        string `json:"teacher"`
	AIContentRatio string `json:"aicontentratio"`
}

func AssessmentSubmissionSubmitACK(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//获取key和任课老师
	var info assessmentACKBody
	if !util.JsonCtx(ctx, &info) {
		util.DebugError("AssessmentSubmissionSubmitACK", "解析post body失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常", nil)
		return
	}
	//ai率不能为空
	if len(info.AIContentRatio) == 0 {
		util.DebugError("AI率不上传也能绕过?", id)
		util.ResponseNAK_MSG(ctx, "请上传AI使用率!", nil)
		return
	}
	//记录进数据库
	indices, err := service.AssessmentSubmitACK(id, info.Key, info.Teacher, info.AIContentRatio)
	if err != nil {
		util.DebugError("AssessmentSubmissionSubmitACK", err)
		util.ResponseNAK_MSG(ctx, "服务器异常", nil)
		return
	}
	//写入成功
	util.ResponseACK_MSG(ctx, "上传成功!", map[string]interface{}{
		"indices": indices,
	})
	util.DebugSuccess(id, "上传考核代码成功!")
}
