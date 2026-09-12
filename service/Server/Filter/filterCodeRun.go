package Filter

import (
	"fmt"
	data "newaoe/service/Data"
	"newaoe/util"

	"github.com/gin-gonic/gin"
)

// 防止用户一直提交代码进行攻击
func FilterCodeRun() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//
		userInfo := util.GetCtxTookenInfo(ctx)
		if userInfo == nil {
			util.ResponseNAK_MSG(ctx, "cookie非法或者错误", "")
			ctx.Abort()
			return
		}
		//
		id := userInfo["id"]
		ttl := 0 //data.CodeRunRecordTTL(id)
		if ttl > 0 {
			util.ResponseNAK_MSG(ctx, fmt.Sprintf("%v同学,你提交的太快了!请在%d秒后再提交!", id, ttl), "")
			ctx.Abort()
			return
		} else {
			//不管之后也没有成功，我们均在这加入记录
			if !data.CodeRunRecordAdd(id) {
				util.ResponseNAK_MSG(ctx, "提交失败,请重试!", "")
				ctx.Abort()
				return
			}
		}
		//
		//
		ctx.Next()
	}
}
