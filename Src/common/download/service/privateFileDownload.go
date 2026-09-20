package service

import (
	"newaoe/Src/config"
	"newaoe/Src/util"
	"path"
	"strings"
	"time"
)

// 下载学生私有资源
func PrivateFileDownload(filepath string, userInfo map[string]interface{}, attachment bool) string {
	//鉴权
	if !privateDownloadCheck(userInfo, filepath) {
		util.DebugError(userInfo["id"].(string), "无权访问该资源!")
		return ""
	}
	//
	return DownloadFile(filepath, time.Duration(config.Conf.Other.PrivateFileDownloadUrlExpireTime)*time.Minute, attachment)
}

// 私人资源下载鉴权(目前只不允许看别人的代码)
func privateDownloadCheck(userInfo map[string]interface{}, filePath string) bool {
	id := userInfo["id"].(string)
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
	return !util.IsPrefix(filePath, codePath)
}
