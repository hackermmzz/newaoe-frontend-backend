package service

import (
	"newaoe/Src/config"
	"newaoe/Src/util"
	"time"
)

// 下载公共资源
func PublicFileDownload(filepath string, userInfo map[string]interface{}, attachment bool) string {
	//鉴权
	if !publicFileDownloadCheck(userInfo, filepath) {
		util.DebugError(userInfo["id"].(string), "无权访问该资源!")
		return ""
	}
	//
	return DownloadFile(filepath, time.Duration(config.Conf.Other.PublicFileDownloadUrlExpireTime)*time.Second, attachment)
}

// 公共资源下载鉴权(目前没有什么鉴权,后续可以添加一些访问频率限制之类的)
func publicFileDownloadCheck(userInfo map[string]interface{}, filePath string) bool {
	//目前公共资源随便下载,后续可以添加一些访问频率限制之类的
	return true
}
