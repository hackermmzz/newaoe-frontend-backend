package Upload

import (
	"newaoe/dao"
	data "newaoe/service/Data"
	"newaoe/util"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func ProcessUploadAvatar(ctx *gin.Context, userInfo map[string]string, file minio.ObjectInfo) {
	id, _ := userInfo["id"]
	//
	var err error
	session := dao.DB.NewSession()
	defer session.Rollback()
	if err = session.Begin(); err != nil {
		util.ResponseNAK_MSG(ctx, "服务器异常", "")
		return
	}
	//写入数据库
	if data.UserResetAvatar(session, id, file.Key) {
		if err = session.Commit(); err != nil {
			util.ResponseNAK_MSG(ctx, "服务器异常", "")
			return
		}
		util.ResponseACK_MSG(ctx, "头像更改成功", "")
	} else {
		util.ResponseNAK_MSG(ctx, "头像上传失败", "")
	}
}

func AvatarUploadConfirm(ctx *gin.Context, userInfo map[string]string, path string) {
	info := dao.OssCheckFileExist(path)
	if info == nil {
		util.ResponseNAK_MSG(ctx, "头像上传确认失败!", "")
		util.Debug("AvatarUploadConfirm:", userInfo["id"], path)
		return
	}
	ProcessUploadAvatar(ctx, userInfo, *info)
}
