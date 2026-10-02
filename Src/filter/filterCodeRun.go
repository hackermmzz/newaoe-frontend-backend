package filter

import (
	"context"
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
		//如果阻塞提交，直接返回
		if blocked, err := queryBlock("BlockCodeRun"); blocked || err != nil {
			if err == nil {
				err = util.NewError("服务器异常!")
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
	ttl, err := redis.RedisTTL(context.Background(), fmt.Sprintf("CodeRunOrSubmit:%v", id))
	if err != nil {
		util.DebugError("CodeRunRecordTTL:", err)
		return int64(1e9)
	}
	return int64(ttl.Seconds())
}

func codeRunRecordAdd(id string) bool {
	ok := redis.RedisSet(context.Background(), fmt.Sprintf("CodeRunOrSubmit:%v", id), "", time.Duration(config.Conf.Code.CodeSubmitInterval)*time.Second)
	//
	return ok
}

func queryBlock(key string) (bool, error) {
	var queryBlockLua = `
	local val = redis.call("GET", KEYS[1])
	if val == nil then
		return 0
	end
	return tonumber(val)
	`
	resCMD := redis.RedisLua(context.Background(), queryBlockLua, []string{key})
	res, err := resCMD.Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}
