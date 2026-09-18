package initialize

import "newaoe/Src/redis"

func RedisInit() {
	//连接Redis缓存
	redis.ConnectRedis()
}
