package controller

import (
	"newaoe/Src/codeAssessmentSubmit/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func AssessmentSubmissionSubmit(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//获取上传链接以及信息key
	urls, key, err := service.AssessmentSubmit(id)
	//回复
	if err != nil {
		util.DebugError("AssessmentSubmissionSubmit失败!", err)
		util.ResponseNAK_MSG(ctx, "服务器异常", nil)
		return
	}
	util.ResponseACK_MSG(ctx, "获取链接成功!",
		map[string]interface{}{
			"urls": urls,
			"key":  key,
		})
}
