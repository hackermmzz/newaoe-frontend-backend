package filter

import (
	"context"
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/redis"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

var limitCodeSubmitLua = `
local current = tonumber(redis.call("GET", KEYS[1]) or "0")
local limit = tonumber(ARGV[1])
local expire = tonumber(ARGV[2])
if current ==nil then
	return {-1, -1}
if current >= limit then
    return {-1, current}
end
local newValue = redis.call("INCR", KEYS[1])
if newValue == 1 then
    redis.call("EXPIRE", KEYS[1], expire)
end

return {1, newValue}
`

// 判断是否达到提交/运行限制
func FilterLimitCodeSubmit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//如果阻塞提交，直接返回
		if blocked, err := queryBlock("BlockCommonSubmit"); blocked || err != nil {
			if err == nil {
				err = util.NewError("禁止提交!")
			}
			util.ResponseNAK_MSG(ctx, err.Error(), nil)
			ctx.Abort()
			return
		}
		//
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
		result, err := redis.RedisLua(
			context.Background(),
			limitCodeSubmitLua,
			[]string{key},
			config.Conf.Code.CodeSubmitTimesPerDay,
			util.GetLeftTimeForOneDay().Seconds(),
		).Result()

		if err != nil {
			util.DebugError("FilterLimitCodeSubmit", err)
			util.ResponseNAK_MSG(ctx, "服务器异常!", err.Error())
			ctx.Abort()
			return
		}

		data := result.([]interface{})
		status := data[0].(int64)

		if status == 1 {
			ctx.Next()
		} else {
			util.ResponseNAK_MSG(ctx, "当天提交次数已耗尽!", nil)
			ctx.Abort()
		}
	}
}
