package Filter

import (
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"

	"github.com/gin-gonic/gin"
)

// 代码运行过滤器，过滤评测提交的代码，防止用户多次上传，提前上传
func FilterAssessmentSubmission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取用户
		info := util.GetCtxTookenInfo(ctx)
		if info == nil {
			util.Debug("为什么这一步会出现这种异常，后端有漏洞!")
			util.ResponseNAK_MSG(ctx, "cookie异常!", "")
			ctx.Abort()
			return
		}
		//获取用户id
		id := info["id"]
		//判断是否到达可以提交评测的时间(默认不加以限制)

		//判断是否已经提交过评测了
		if len(dao.CodeAssessmentGetByID(id)) >= config.Conf.Code.CodeAssessmentTimes {
			util.ResponseNAK_MSG(ctx, "你已经提交过评测了哦!", "")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
