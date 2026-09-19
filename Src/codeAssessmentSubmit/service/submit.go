package service

import (
	"errors"
	"newaoe/Src/common/upload"
	"newaoe/Src/config"
	"newaoe/Src/util"
	"path"
	"time"
)

// 返回生成的URL和key
func AssessmentSubmit(id string) ([]string, string, error) {
	prefix := "AssessmentSubmission" + "_" + util.UTC_Time().Format("2006_01_02_150405000")
	base_dir := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCodeFolder, prefix)
	//
	header := path.Join(base_dir, "mmzz.h")
	source := path.Join(base_dir, "mmzz.cpp")
	//
	duration := time.Duration(config.Conf.Other.CodeCommonSubmitUrlExpireTime) * time.Second
	urls, key := upload.UploadFileWithCachePath(id, []string{header, source}, duration)
	if len(urls) != 2 {
		return nil, "", errors.New("返回链接失败!")
	}
	//
	return urls, key, nil
}
