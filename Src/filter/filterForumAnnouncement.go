package filter

import (
	"context"
	"fmt"
	"newaoe/Src/redis"
	"newaoe/Src/user/dao"
	"newaoe/Src/user/model"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

// key 已存在时返回剩余 TTL（秒）。
// key 不存在时先写入 value 0（过期时间取 ARGV[1]，单位秒），再返回 0。
// 脚本遇到 return 会结束，所以必须先创建再返回。
var forumAnnouncementLua = `
if redis.call("EXISTS", KEYS[1]) == 1 then
	return redis.call("TTL", KEYS[1])
end
redis.call("SET", KEYS[1], 0, "EX", ARGV[1])
return 0
`

func FilterForumAnnouncement() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//限制帖子发送频率(获取ID)
		userInfo := util.GetCtxTookenInfo(ctx)
		if userInfo == nil {
			util.ResponseNAK_MSG(ctx, "cookie非法或者错误", "")
			ctx.Abort()
			return
		}
		id := userInfo["id"].(string)
		//判断是否是VIP用户，如果是则不拦截
		user, err := dao.UserGet(nil, id)
		if err == nil && user != nil && user.Vip >= model.VIP_SUPER {
			ctx.Next()
			return
		}
		//
		key := fmt.Sprintf("forum_announcement:%s", id)
		expireSeconds := int64(3 * 60)
		ttl, err := redis.RedisLua(
			context.Background(),
			forumAnnouncementLua,
			[]string{key},
			expireSeconds, //3分钟
		).Int64()
		if err != nil {
			util.DebugError("FilterForumAnnouncement", err)
			util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
			ctx.Abort()
			return
		}
		//返回 0 表示 key 原本不存在，本次已创建，允许发帖。
		//返回剩余 TTL 表示仍在冷却中。
		if ttl != 0 {
			util.ResponseNAK_MSG(ctx, fmt.Sprintf("%v同学,你发帖太快了!请在%d秒后再发!", id, ttl), nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
