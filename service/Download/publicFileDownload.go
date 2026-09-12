package Download

import (
	"newaoe/config"
	"newaoe/util"
	"time"

	"github.com/gin-gonic/gin"
)

// 下载公共资源
func PublicFileDownload(ctx *gin.Context, filepath string, uerInfo map[string]string, attachment bool) {
	//鉴权
	if !PublicFileDownloadCheck(uerInfo, filepath) {
		util.ResponseNAK_MSG(ctx, "无权访问该资源!", "")
		return
	}
	//
	DownloadFile(ctx, filepath, time.Duration(config.Conf.Other.PublicFileDownloadUrlExpireTime)*time.Second, attachment)
}

// 公共资源下载鉴权(目前没有什么鉴权,后续可以添加一些访问频率限制之类的)
func PublicFileDownloadCheck(userInfo map[string]string, filePath string) bool {
	//目前公共资源随便下载,后续可以添加一些访问频率限制之类的
	return true
}
