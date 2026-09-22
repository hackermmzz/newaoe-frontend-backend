package service

import (
	"newaoe/Src/config"
	"newaoe/Src/util"
)

// 获取路径的类别category
func GetFilePathCategory(path string) string {
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
