package filter

import (
	"newaoe/Src/codeAssessmentSubmit/dao"
	"newaoe/Src/config"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

// 代码运行过滤器，过滤评测提交的代码，防止用户多次上传，提前上传
func FilterLimitAssessmentSubmission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//如果阻塞提交，直接返回
		if blocked, err := queryBlock("BlockAssessmentSubmit"); blocked || err != nil {
			if err == nil {
				err = util.NewError("禁止提交!")
			}
			util.ResponseNAK_MSG(ctx, err.Error(), nil)
			ctx.Abort()
			return
		}
		//获取用户
		info := util.GetCtxTookenInfo(ctx)
		if info == nil {
			util.DebugError("为什么这一步会出现这种异常，后端有漏洞!")
			util.ResponseNAK_MSG(ctx, "cookie异常!", "")
			ctx.Abort()
			return
		}
		//获取用户id
		id := info["id"].(string)
		//判断是否到达可以提交评测的时间(默认不加以限制)

		//判断是否已经提交过评测了
		records, err := dao.CodeAssessmentGetByID(nil, id)
		if err != nil {
			util.ResponseNAK_MSG(ctx, err.Error(), nil)
			ctx.Abort()
			return
		}
		if len(records) >= config.Conf.Code.CodeAssessmentTimes {
			util.ResponseNAK_MSG(ctx, "你已经提交过评测了哦!", "")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
