package filter

import (
	"newaoe/Src/user/model"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

// 过滤掉权限小于等于游客的用户
func FilterTourist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//
		userInfo := util.GetCtxTookenInfo(ctx)
		if userInfo == nil {
			util.ResponseNAK_MSG(ctx, "cookie非法或者错误", "")
			ctx.Abort()
			return
		}
		//
		vip := int(userInfo["vip"].(float64))
		if vip <= model.VIP_TOURIST {
			util.ResponseNAK_MSG(ctx, "游客等级以下用户无此业务权限!", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
