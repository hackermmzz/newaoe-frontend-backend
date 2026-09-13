package Upload

import (
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func FileUpload(ctx *gin.Context) {
	//获取类别
	pathParts := strings.Split(ctx.Request.URL.Path, "/")
	category := pathParts[len(pathParts)-1]
	//获取用户信息
	userInfo := util.GetCtxTookenInfo(ctx)
	if userInfo == nil {
		util.ResponseNAK_MSG(ctx, "cookie非法或者错误", "")
		return
	}
	//
	//
	switch category {
	case config.Conf.User.UserAvatarFolder:
		AvatarUpload(ctx, userInfo)
		return
	case config.Conf.User.UserCodeFolder:
		CodeUpload(ctx, userInfo)
		return
	case config.Conf.Other.AssessmentCodeUploadCategory:
		AssessmentSubmissionUpload(ctx, userInfo)
		return
	default:
		//
		util.ResponseNAK_MSG(ctx, "上传文件失败", "")
	}

}

// 文件上传接口(返回上传链接)
func UploadFile(ctx *gin.Context, filePath []string, expireDuration []time.Duration) bool {
	urls := dao.GetUploadFileUrls(filePath, expireDuration)
	if urls == nil {
		util.Debug("UploadFile:生成上传链接失败")
		util.ResponseNAK_MSG(ctx, "上传文件失败!", "")
		return false
	}
	util.ResponseACK_MSG(ctx, "上传文件成功!", map[string]interface{}{"urls": urls})
	return true
}
