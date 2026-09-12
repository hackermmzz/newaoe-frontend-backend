package Download

import (
	"newaoe/config"
	"newaoe/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 下载学生私有资源
func PrivateFileDownload(ctx *gin.Context, filepath string, uerInfo map[string]string, attachment bool) {
	//鉴权
	if !PrivateDownloadCheck(uerInfo, config.Conf.OSS.PrivateBaseFolder, filepath) {
		util.ResponseNAK_MSG(ctx, "无权访问该资源!", "")
		return
	}
	//
	DownloadFile(ctx, filepath, time.Duration(config.Conf.Other.PrivateFileDownloadUrlExpireTime)*time.Second, attachment)
}

// 私人资源下载鉴权(目前只能访问自己的私人资源,后续可以添加一些访问频率限制之类的)
func PrivateDownloadCheck(userInfo map[string]string, category string, filePath string) bool {
	id := userInfo["id"]
	//判断是否是自己的资源
	parts := strings.Split(filePath, "/")
	if len(parts) < 3 {
		return false
	}
	if parts[1] != id {
		return false
	}
	return true
}
