package filter

import (
	UserModel "newaoe/Src/user/model"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func FilterVIP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dt := util.GetCtxTookenInfo(ctx)
		//判断是否是VIP
		vip := int(dt["vip"].(float64))
		if vip < UserModel.VIP_SUPER {
			util.DebugError("有一个不是vip的访问了!", dt["id"].(string))
			util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
