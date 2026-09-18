package controller

import (
	"newaoe/Src/common/upload"
	"newaoe/Src/config"
	"newaoe/Src/util"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func AssessmentSubmissionSubmit(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//
	prefix := "AssessmentSubmission" + "_" + util.UTC_Time().Format("2006_01_02_150405000")
	base_dir := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCodeFolder, prefix)
	//
	header := path.Join(base_dir, "mmzz.h")
	source := path.Join(base_dir, "mmzz.cpp")
	//
	duration := time.Duration(config.Conf.Other.CodeCommonSubmitUrlExpireTime) * time.Second
	urls, key := upload.UploadFileWithCachePath(id, []string{header, source}, duration)
	if len(urls) != 2 {
		util.DebugError("AssessmentSubmissionSubmit 返回链接失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常", nil)
		return
	}
	//回复
	util.ResponseACK_MSG(ctx, "获取链接成功!", map[string]interface{}{"urls": urls, "key": key})
}
