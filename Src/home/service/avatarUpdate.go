package service

import (
	"fmt"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/oss"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
	"path"
	"time"
)

func AvatarUpdate(id string) (string, error) {
	avatarname := fmt.Sprintf("avatar_%s.png", util.GetTimeNanoStr())
	avatarNewPath := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserAvatarFolder, avatarname)
	url := oss.GetUploadFileUrl(avatarNewPath, time.Duration(30)*time.Minute) //30分钟过期
	if url == "" {
		return "", util.NewError("获取上传头像链接失败!")
	}
	//
	var err error
	session := database.NewSession()
	defer session.Close()
	if err = session.Begin(); err != nil {
		return "", util.NewError("服务器异常!")
	}
	defer session.Rollback()
	//写入数据库
	if dao.UserResetAvatar(session, id, avatarNewPath) {
		if err = session.Commit(); err != nil {
			return "", util.NewError("服务器异常!")
		}
	} else {
		return "", util.NewError("服务器异常!")
	}
	//
	return url, nil
}
