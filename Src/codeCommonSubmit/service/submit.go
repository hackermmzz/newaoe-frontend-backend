package service

import (
	"errors"
	"fmt"
	"newaoe/Src/common/upload"
	"newaoe/Src/config"
	"newaoe/Src/oss"
	"newaoe/Src/util"
	"path"
	"time"
)

func CodeCommonSubmit(id string) ([]*oss.UploadPolicyResult, string, error) {
	prefix := fmt.Sprintf("CommonSubmit_%s", util.UUID())
	base_dir := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCodeFolder, prefix)
	header := path.Join(base_dir, "mmzz.h")
	source := path.Join(base_dir, "mmzz.cpp")
	description := path.Join(base_dir, "mmzz.txt")
	//返回信息
	duration := time.Duration(config.Conf.Other.CodeCommonSubmitUrlExpireTime) * time.Second
	urls, key := upload.UploadFileWithCachePath(id, []string{header, source, description}, duration, 0, config.Conf.Server.MaxFilsSizeUserUpload)
	if len(urls) != 3 {
		return nil, "", errors.New("获取上传链接失败!")
	}
	return urls, key, nil
}
