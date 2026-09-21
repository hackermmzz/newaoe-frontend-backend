package service

import (
	"net/url"
	"newaoe/Src/config"
	"newaoe/Src/oss"
	"newaoe/Src/util"
	"time"
)

var (
	DownloadURLType_UnKnown = 0 //未知
	DownloadURLType_OSS     = 1 //oss下载
	DownloadURLType_CDN     = 2 //cdn下载
)

// 用户下载文件接口(返回下载链接),如果attachment为true表示需要当作附件下载
func DownloadFile(filePath string, expireDuration time.Duration, attachment bool) (string, int) {
	if config.Conf.CDN.CDNSwitch {
		//走cdn
		url, err := url.JoinPath(config.Conf.CDN.Host, filePath)
		if err != nil {
			util.DebugError("DownloadFile", err)
			return "", DownloadURLType_UnKnown
		}
		return url, DownloadURLType_CDN
	} else {
		//直接走oss
		url := oss.OssGetDownloadFileUrl(filePath, expireDuration, attachment)
		if url == "" {
			util.DebugError("DownloadFile:无法获取下载链接!")
			return "", DownloadURLType_UnKnown
		}
		return url, DownloadURLType_OSS
	}
}
