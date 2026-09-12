package Upload

import (
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func AssessmentSubmissionUpload(ctx *gin.Context, userInfo map[string]string) {
	id, _ := userInfo["id"]
	//限制提交次数
	if AssessmentSubmissionLimit(id) {
		util.ResponseNAK_MSG(ctx, "你已提交，不可重复提交!", nil)
		return
	}
	//
	prefix := "AssessmentSubmission" + "_" + util.UTC_Time().Format("2006_01_02_150405000")
	base_dir := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCodeFolder, prefix)
	//
	header := path.Join(base_dir, "mmzz.h")
	source := path.Join(base_dir, "mmzz.cpp")
	duration := time.Duration(config.Conf.Other.AssessmentCodeUploadUrlExpireTime) * time.Second
	UploadFile(ctx, []string{header, source}, []time.Duration{duration, duration})
}

func AssessmentSubmissionLimit(id string) bool {
	ret := dao.CodeAssessmentGetByID(id)
	return len(ret) >= config.Conf.Code.CodeAssessmentTimes
}
