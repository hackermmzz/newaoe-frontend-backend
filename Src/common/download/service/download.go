package service

import (
	"newaoe/Src/oss"
	"newaoe/Src/util"
	"time"
)

// 用户下载文件接口(返回下载链接),如果attachment为true表示需要当作附件下载
func DownloadFile(filePath string, expireDuration time.Duration, attachment bool) string {
	url := oss.OssGetDownloadFileUrl(filePath, expireDuration, attachment)
	if url == "" {
		util.DebugError("DownloadFile:无法获取下载链接!")
		return ""
	}
	return url
}
