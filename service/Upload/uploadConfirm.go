package Upload

import (
	"newaoe/config"
	"newaoe/util"

	"github.com/gin-gonic/gin"
)

func FileUploadConfirm(ctx *gin.Context) {
	//获取类别
	category := ctx.Param("category")
	//获取用户信息
	userInfo := util.GetCtxTookenInfo(ctx)
	if userInfo == nil {
		util.ResponseNAK_MSG(ctx, "cookie非法或者错误", "")
		return
	}
	//获取上传的url
	type Msg struct {
		URLs []string `json:"urls"`
	}
	var msg Msg
	if !util.JsonCtx(ctx, &msg) {
		util.ResponseNAK_MSG(ctx, "请求体格式错误!", "")
		return
	}
	//解析url获取路径
	var path []string
	for _, p := range msg.URLs {
		p = util.GetUrlFilePath(p) //注意这里拆分出来还有一个bucket
		p = util.SplitBucketAndFilePath(config.Conf.OSS.BucketName, p)
		if p == "" {
			util.ResponseNAK_MSG(ctx, "保存的路径错误!", "")
			return
		}
		path = append(path, p)
	}
	//要保证传入的路径数量与对应的业务需要数量一致
	switch category {
	////////////////
	case config.Conf.User.UserAvatarFolder:
		if len(path) != 1 {
			util.ResponseNAK_MSG(ctx, "你小子爬虫我吗!", "")
			return
		}
		AvatarUploadConfirm(ctx, userInfo, path[0])
		return
		////////////////
	case config.Conf.User.UserCodeFolder:
		if len(path) != 3 {
			util.ResponseNAK_MSG(ctx, "你小子爬虫我吗!", "")
			return
		}
		CodeUploadConfirm(ctx, userInfo, path)
		return
	case config.Conf.Other.AssessmentCodeUploadCategory:
		if len(path) != 2 {
			util.ResponseNAK_MSG(ctx, "你小子爬虫我吗!", "")
			return
		}
		AssessmentSubmissionUploadConfirm(ctx, userInfo, path)
	default:
		util.ResponseNAK_MSG(ctx, "上传文件确认失败", "")
	}
}
