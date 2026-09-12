package Upload

import (
	"newaoe/config"
	"newaoe/service/Code"
	"newaoe/util"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func CodeUpload(ctx *gin.Context, userInfo map[string]string) {

	id, _ := userInfo["id"]
	//限制提交次数
	if Code.LimitCodeUploadOrRun(id) {
		util.ResponseNAK_MSG(ctx, "当天3次提交次数已消耗尽!", nil)
		return
	}
	//
	prefix := util.UTC_Time().Format("2006_01_02_150405000")
	base_dir := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCodeFolder, prefix)
	//
	header := path.Join(base_dir, "mmzz.h")
	source := path.Join(base_dir, "mmzz.cpp")
	description := path.Join(base_dir, "mmzz.txt")
	duration := time.Duration(config.Conf.Other.CodeCommonUploadUrlExpireTime) * time.Second
	UploadFile(ctx, []string{header, source, description}, []time.Duration{duration, duration, duration})
}
