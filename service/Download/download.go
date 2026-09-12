package Download

import (
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

// 下载文件路由
func FileDownload(ctx *gin.Context) {
	category := ctx.Param("category")
	filepathParam := ctx.Param("filepath")
	attachment := ctx.Query("download") == "true"
	path := path.Join(category, filepathParam)
	//检查路径是否合法
	userInfo := util.GetCtxTookenInfo(ctx)
	if userInfo == nil {
		util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
		return
	}
	//资源直接下载
	switch category {
	//公共文件下载
	case config.Conf.OSS.PublicBaseFolder:
		PublicFileDownload(ctx, path, userInfo, attachment)
		return
	//私人文件下载
	case config.Conf.OSS.PrivateBaseFolder:
		PrivateFileDownload(ctx, path, userInfo, attachment)
		return
	default:
		util.ResponseNAK_MSG(ctx, "不支持的下载类别!", "居然戏弄我!")
	}
}

// 判断路径是否合法
func UserDownloadCheck(userInfo map[string]string, category string, filePath string) bool {
	/*
		//
		basedir := ""
		targetdir := ""
		//检查权限(如果不是公共路径,那么不能相互之间访问)
		switch category {
		case config.Conf.OSS.PublicBaseFolder:
			//公共目录随便下载
			basedir = path.Join(config.Conf.Data.DataBaseFolder, config.Conf.Data.PublicBaseFolder)
			targetdir = path.Join(basedir, filePath)
		case config.Conf.Data.PrivateBaseFolder:
			//获取用户
			id := userInfo["id"]
			basedir = path.Join(config.Conf.Data.DataBaseFolder, config.Conf.Data.PrivateBaseFolder, id) //保证一定是这个用户的私人文件夹
			targetdir = path.Join(config.Conf.Data.DataBaseFolder, config.Conf.Data.PrivateBaseFolder, filePath)

		default:
			return false
		}
		//判断是否为basedir的子目录
		if !util.IsSubDir(basedir, targetdir) {
			return false
		}
		//判断文件是否存在
		return util.IsFileExist(targetdir)
	*/
	return true
}

// 用户下载文件接口(返回下载链接),如果attachment为true表示需要当作附件下载
func DownloadFile(ctx *gin.Context, filePath string, expireDuration time.Duration, attachment bool) {
	url := dao.OssGetDownloadFileUrl(filePath, expireDuration, attachment)
	if url == "" {
		util.Debug("DownloadFile:无法获取下载链接!")
		util.ResponseNAK_MSG(ctx, "下载链接获取错误!", "")
		return
	}
	util.ResponseACK_MSG(ctx, "请求成功!", map[string]interface{}{"url": url})
}
