package service

import (
	"newaoe/Src/config"
	"newaoe/Src/user/model"
	"newaoe/Src/util"
	"path"
	"strings"
)

func DownloadAuth(id string, vip int, filepath string) (bool, error) {
	if isVIP(vip) { //超级用户随便下载
		return true, nil
	} else if isPublicPath(filepath) { //公开路径随便下载
		return true, nil
	} else if isTempPath(filepath) { //临时文件除了超级管理员可以下载，其他人都不允许下载
		return false, nil
	} else if isCodePath(filepath) { //代码除本人以外不允许下载
		return util.IsPrefix(filepath, path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCodeFolder)), nil
	} else {
		//随便下载
		return true, nil
	}
}

// 判断vip等级是否足够
func isVIP(vip int) bool {
	return vip >= model.VIP_SUPER
}

// 判断路径是否是代码路径
func isCodePath(filepath string) bool {
	arr := strings.Split(filepath, "/")
	if len(arr) < 3 {
		return false
	}
	return arr[2] == config.Conf.User.UserCodeFolder
}

// 判断是否是公开路径
func isPublicPath(filepath string) bool {
	return util.IsPrefix(filepath, config.Conf.OSS.PublicBaseFolder)
}

// 判断是否是临时文件路径
func isTempPath(filepath string) bool {
	return util.IsPrefix(filepath, config.Conf.OSS.TmpBaseFolder)
}
