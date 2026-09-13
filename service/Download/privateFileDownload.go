package Download

import (
	"newaoe/config"
	"newaoe/util"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 下载学生私有资源
func PrivateFileDownload(ctx *gin.Context, filepath string, uerInfo map[string]string, attachment bool) {
	//鉴权
	if !PrivateDownloadCheck(uerInfo, filepath) {
		util.ResponseNAK_MSG(ctx, "无权访问该资源!", "")
		return
	}
	//
	DownloadFile(ctx, filepath, time.Duration(config.Conf.Other.PrivateFileDownloadUrlExpireTime)*time.Second, attachment)
}

// 私人资源下载鉴权(目前只不允许看别人的代码)
func PrivateDownloadCheck(userInfo map[string]string, filePath string) bool {
	id := userInfo["id"]
	//判断是否是自己的资源
	parts := strings.Split(filePath, "/")
	if len(parts) < 3 {
		return false
	}
	if parts[1] == id {
		return true
	}
	//否则判断是否是代码
	codePath := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCodeFolder)
	return !strings.HasPrefix(filePath, codePath)
}
