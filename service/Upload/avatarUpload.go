package Upload

import (
	"newaoe/config"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func AvatarUpload(ctx *gin.Context, userInfo map[string]string) {
	id := userInfo["id"]
	//覆盖原来的头像
	targetdir := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserAvatarFolder, "avatar.png")
	duration := time.Duration(config.Conf.Other.AvatarUploadUrlExpireTime) * time.Second
	UploadFile(ctx, []string{targetdir}, []time.Duration{duration})
}
