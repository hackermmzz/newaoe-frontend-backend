package Filter

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"

	"github.com/gin-gonic/gin"
)

// 判断是否达到提交/运行限制
func FilterLimitCodeUploadOrRun() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInfo := util.GetCtxTookenInfo(ctx)
		if userInfo == nil {
			util.ResponseNAK_MSG(ctx, "cookie非法或者错误", "")
			ctx.Abort()
			return
		}
		id := userInfo["id"]
		//超级用户不管
		for _, superuser := range config.Conf.Other.SuperUser {
			if superuser == id {
				ctx.Next()
				return
			}
		}
		//
		key := fmt.Sprintf("CommonUploadOrRunTimes_%v", id)
		time, exist := dao.RedisGet(context.Background(), key)
		curSubmitTime := 0
		if exist {
			err := json.Unmarshal(time, &curSubmitTime)
			if err != nil {
				util.Debug("LimitCodeUploadOrRun Unmarshal Error!", err)
			}
		}
		//判断是否达到限制
		if curSubmitTime >= config.Conf.Code.CodeSubmitTimesPerDay {
			util.ResponseNAK_MSG(ctx, "你已经达到当天提交限制!", nil)
			ctx.Abort()
			return
		}
		//
		dao.RedisSet(context.Background(), key, curSubmitTime+1, util.GetLeftTimeForOneDay())
		ctx.Next()
	}
}
