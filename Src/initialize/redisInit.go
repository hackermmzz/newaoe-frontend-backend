package initialize

import (
	"newaoe/Src/redis"
	"newaoe/Src/util"
)

func RedisInit() {
	//连接Redis缓存
	if err := redis.ConnectRedis(); err != nil {
		util.DebugError("RedisInit:", err)
		panic(err)
	}
}
