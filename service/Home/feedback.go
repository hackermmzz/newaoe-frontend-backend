package Home

import (
	"newaoe/config"
	"newaoe/service/Upload"
	"newaoe/util"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func StudentFeedback(ctx *gin.Context) {
	//获取用户id
	userInfo := util.GetCtxTookenInfo(ctx)
	if userInfo == nil {
		util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
		return
	}
	id := userInfo["id"]
	//指定文件名称
	filename := util.UTC_Time().Format("20060102_150405") + "_" + id + ".json"
	targetfile := path.Join(config.Conf.OSS.PublicBaseFolder, "feedback", filename)
	//走upload的路线,但是不会confirm(url过期时间和普通代码提交时间一样)
	duration := time.Duration(config.Conf.Other.CodeCommonUploadUrlExpireTime) * time.Second
	Upload.UploadFile(ctx, []string{targetfile}, []time.Duration{duration})
}
