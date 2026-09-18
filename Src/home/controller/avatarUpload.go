package controller

import (
	"fmt"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/oss"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func AvatarUpdate(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id, _ := userInfo["id"].(string)
	//生成上传链接
	avatarname := fmt.Sprintf("avatar_%.png", util.GetTimeNanoStr())
	avatarNewPath := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserAvatarFolder, avatarname)
	url := oss.GetUploadFileUrl(avatarNewPath, time.Duration(30)*time.Minute) //30分钟过期
	if url == "" {
		util.DebugError("获取上传头像链接失败!")
		util.ResponseNAK_MSG(ctx, "获取上传链接失败!", nil)
		return
	}
	//
	var err error
	session := database.NewSession()
	defer session.Rollback()
	if err = session.Begin(); err != nil {
		util.ResponseNAK_MSG(ctx, "服务器异常", "")
		return
	}
	//写入数据库
	if dao.UserResetAvatar(session, id, avatarNewPath) {
		if err = session.Commit(); err != nil {
			util.ResponseNAK_MSG(ctx, "服务器异常", "")
			return
		}
		util.ResponseACK_MSG(ctx, "头像更改成功", "")
	} else {
		util.ResponseNAK_MSG(ctx, "头像上传失败", "")
	}
}
