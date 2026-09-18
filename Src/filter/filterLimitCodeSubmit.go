package filter

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/redis"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

// 判断是否达到提交/运行限制
func FilterLimitCodeSubmit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userInfo := util.GetCtxTookenInfo(ctx)
		if userInfo == nil {
			util.ResponseNAK_MSG(ctx, "cookie非法或者错误", "")
			ctx.Abort()
			return
		}
		id := userInfo["id"].(string)
		//超级用户不管
		_, ok1 := config.SuperUserMap[id]
		_, ok2 := config.UnlimitUserMap[id]
		if ok1 || ok2 {
			ctx.Next()
			return
		}
		//
		key := fmt.Sprintf("CommonUploadOrRunTimes_%v", id)
		time, exist := redis.RedisGet(context.Background(), key)
		curSubmitTime := 0
		if exist {
			err := json.Unmarshal(time, &curSubmitTime)
			if err != nil {
				util.DebugError("LimitCodeUploadOrRun Unmarshal Error!", err)
			}
		}
		//判断是否达到限制
		if curSubmitTime >= config.Conf.Code.CodeSubmitTimesPerDay {
			util.ResponseNAK_MSG(ctx, "你已经达到当天提交限制!", nil)
			ctx.Abort()
			return
		}
		//
		redis.RedisSet(context.Background(), key, curSubmitTime+1, util.GetLeftTimeForOneDay())
		ctx.Next()
	}
}
