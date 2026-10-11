package service

import (
	"newaoe/Src/config"
	"time"
)

// 下载公共资源
func PublicFileDownload(filepath string, id string, attachment bool) (string, int) {
	return DownloadFile(filepath, time.Duration(config.Conf.Other.PublicFileDownloadUrlExpireTime)*time.Minute, attachment)
}
