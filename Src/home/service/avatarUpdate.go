package service

import (
	"fmt"
	"newaoe/Src/common/upload"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/oss"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
	"path"
	"time"
)

func AvatarUpdate(id string) (*oss.UploadPolicyResult, error) {
	avatarname := fmt.Sprintf("avatar_%s.png", util.UUID())
	avatarNewPath := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserAvatarFolder, avatarname)
	url := upload.UploadPostPolicyFile(avatarNewPath, time.Duration(30)*time.Minute, 0, config.Conf.Server.MaxFilsSizeUserUpload) //30分钟过期
	if url == nil {
		return nil, util.NewError("获取上传头像凭证失败!")
	}
	//
	var err error
	session := database.NewSession()
	defer session.Close()
	if err = session.Begin(); err != nil {
		return nil, util.NewError("服务器异常!")
	}
	defer session.Rollback()
	//写入数据库
	updated, err := dao.UserResetAvatar(session, id, avatarNewPath)
	if err != nil {
		return nil, err
	}
	if updated {
		if err = session.Commit(); err != nil {
			return nil, util.NewError("服务器异常!")
		}
	} else {
		return nil, util.NewError("服务器异常!")
	}
	//
	return url, nil
}
