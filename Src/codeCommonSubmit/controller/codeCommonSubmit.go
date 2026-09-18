package controller

import (
	"newaoe/Src/common/upload"
	"newaoe/Src/config"
	"newaoe/Src/util"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func CodeCommonSubmit(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//生成3个链接分别是header source 以及description
	prefix := util.UTC_Time().Format("2006_01_02_150405000")
	base_dir := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCodeFolder, prefix)
	header := path.Join(base_dir, "mmzz.h")
	source := path.Join(base_dir, "mmzz.cpp")
	description := path.Join(base_dir, "mmzz.txt")
	//返回信息
	duration := time.Duration(config.Conf.Other.CodeCommonSubmitUrlExpireTime) * time.Second
	urls, key := upload.UploadFileWithCachePath(id, []string{header, source, description}, duration)
	if len(urls) != 3 {
		util.DebugError("CodeCommonSubmit 返回链接失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常", nil)
		return
	}
	//回复
	util.ResponseACK_MSG(ctx, "获取链接成功!", map[string]interface{}{"urls": urls, "key": key})
}
