package filter

import (
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"

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
		id := userInfo["id"].(string)
		//如果是超级用户，不管
		_, exist := config.SuperUserMap[id]
		if exist {
			ctx.Next()
			return
		}
		//
		ttl := codeRunRecordTTL(id)
		if ttl > 0 {
			util.ResponseNAK_MSG(ctx, fmt.Sprintf("%v同学,你提交的太快了!请在%d秒后再提交!", id, ttl), "")
			ctx.Abort()
			return
		} else {
			//不管之后也没有成功，我们均在这加入记录
			if !codeRunRecordAdd(id) {
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

func codeRunRecordTTL(id string) int64 {
	ttl, err := redis.RDB.TTL(redis.RDB.Context(), fmt.Sprintf("CodeRunOrSubmit:%v", id)).Result()
	if err != nil {
		util.DebugError("CodeRunRecordTTL:", err)
		return int64(1e9)
	}
	return int64(ttl.Seconds())
}

func codeRunRecordAdd(id string) bool {
	err := redis.RDB.Set(redis.RDB.Context(), fmt.Sprintf("CodeRunOrSubmit:%v", id), "", time.Duration(config.Conf.Code.CodeSubmitInterval)*time.Second).Err()
	//
	if err != nil {
		return false
	}
	return true
}
