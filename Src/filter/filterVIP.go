package filter

import (
	"newaoe/Src/user/model"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func FilterVIP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dt := util.GetCtxTookenInfo(ctx)
		//判断是否是VIP
		id := dt["id"].(string)
		vip := int(dt["vip"].(float64))
		if vip < model.VIP_SUPER {
			util.DebugError("有一个不是vip的访问了!", id)
			util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
