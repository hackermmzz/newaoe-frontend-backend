package service

import (
	"errors"
	"newaoe/Src/common/upload"
	"newaoe/Src/config"
	"newaoe/Src/util"
	"path"
	"time"
)

func CodeCommonSubmit(id string) ([]string, string, error) {
	prefix := util.UTC_Time().Format("2006_01_02_150405000")
	base_dir := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCodeFolder, prefix)
	header := path.Join(base_dir, "mmzz.h")
	source := path.Join(base_dir, "mmzz.cpp")
	description := path.Join(base_dir, "mmzz.txt")
	//返回信息
	duration := time.Duration(config.Conf.Other.CodeCommonSubmitUrlExpireTime) * time.Second
	urls, key := upload.UploadFileWithCachePath(id, []string{header, source, description}, duration)
	if len(urls) != 3 {
		return nil, "", errors.New("获取上传链接失败!")
	}
	return urls, key, nil
}
