package Upload

import (
	"newaoe/config"
	"newaoe/util"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func CodeUpload(ctx *gin.Context, userInfo map[string]string) {

	id, _ := userInfo["id"]
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
