package service

import (
	"newaoe/Src/config"
	"time"
)

// 下载学生私有资源
func PrivateFileDownload(filepath string, id string, attachment bool) (string, int) {
	//
	return DownloadFile(filepath, time.Duration(config.Conf.Other.PrivateFileDownloadUrlExpireTime)*time.Minute, attachment)
}
