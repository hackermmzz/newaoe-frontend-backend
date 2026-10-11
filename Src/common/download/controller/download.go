package controller

import (
	"errors"
	"newaoe/Src/common/download/service"
	"newaoe/Src/config"
	"newaoe/Src/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// 下载文件路由
func FileDownload(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	//获取路径以及路径类别
	filePath := strings.TrimLeft(ctx.Param("filepath"), "/")
	attachment := ctx.Query("attachment") == "true"
	category := getFilePathCategory(filePath)
	if category == "" {
		util.DebugError("category未知!", category, filePath)
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//鉴权
	id := userInfo["id"].(string)
	vip := int(userInfo["vip"].(float64))
	ok, err := service.DownloadAuth(id, vip, filePath)
	if !ok || err != nil {
		if err == nil {
			err = errors.New("下载权限不足!")
		}
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//资源下载
	url := ""
	download_type := service.DownloadURLType_UnKnown
	switch category {
	//公共文件下载
	case config.Conf.OSS.PublicBaseFolder:
		url, download_type = service.PublicFileDownload(filePath, id, attachment)
	//私人文件下载
	case config.Conf.OSS.PrivateBaseFolder:
		url, download_type = service.PrivateFileDownload(filePath, id, attachment)
	//其他文件
	default:
		url, download_type = service.PublicFileDownload(filePath, id, attachment)
	}
	//返回下载链接
	if url == "" {
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	util.ResponseACK_MSG(ctx, "获取下载链接成功!", map[string]interface{}{
		"url":           url,
		"download_type": download_type,
	})
}

// 获取路径的类别category
func getFilePathCategory(path string) string {
	categories := []string{
		config.Conf.OSS.PublicBaseFolder,
		config.Conf.OSS.PrivateBaseFolder,
	}
	for _, category := range categories {
		if util.IsPrefix(path, category) {
			return category
		}
	}
	return ""
}
