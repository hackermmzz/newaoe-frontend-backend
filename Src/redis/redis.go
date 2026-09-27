package redis

import (
	"context"
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/util"

	"time"

	"github.com/go-redis/redis/v8"
)

// 类型定义
type Pipeliner = redis.Pipeliner

// 全局 Redis 客户端（单例，线程安全）
var redisRDB *redis.Client

// ConnectRedis 生产级 Redis 初始化（连接池 + 重试 + 健康检查）
func ConnectRedis() {
	conf := config.Conf.Redis

	// 核心：连接池配置（大厂标准）
	opt := &redis.Options{
		Addr:         conf.Host,
		Password:     conf.Password,
		DB:           0,
		DialTimeout:  500 * time.Millisecond,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		// ===== 连接池核心=====
		PoolSize:     20,              // 最大连接数
		MinIdleConns: 5,               // 最小空闲连接（保持热连接）
		PoolTimeout:  3 * time.Second, // 获取连接超时
		IdleTimeout:  5 * time.Minute, // 空闲连接超时关闭
		MaxRetries:   3,
	}

	// 创建客户端
	redisRDB = redis.NewClient(opt)

	// 健康检查 + 启动重试
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	for i := 0; i < 10; i++ {
		if err = redisRDB.Ping(ctx).Err(); err == nil {
			break
		}
		util.DebugError(fmt.Sprintf("Redis 重试连接 %d/10: %v", i+1, err))
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		panic(fmt.Sprintf("Redis 连接最终失败: %v", err))
	}

	util.DebugSuccess("Redis连接池初始化成功")
}

// 获取上下文
func NewContext() context.Context {
	return redisRDB.Context()
}

// 生成pipeline
func NewTxPipeline() redis.Pipeliner {
	return redisRDB.TxPipeline()
}

// 获取key的ttl
func RedisTTL(ctx context.Context, key string) (time.Duration, error) {
	return redisRDB.TTL(ctx, key).Result()
}

// 原子操作
func RedisTx(ctx context.Context, fn func(Pipeliner) error) error {
	_, err := redisRDB.TxPipelined(ctx, func(pipe Pipeliner) error {
		return fn(pipe)
	})
	return err
}

// lua脚本
func RedisLua(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return redisRDB.Eval(ctx, script, keys, args...)
}

// RedisGet 带 ctx 规范获取（支持链路超时、链路追踪）
func RedisGet(ctx context.Context, key string) ([]byte, bool) {
	val, err := redisRDB.Get(ctx, key).Bytes()

	// key不存在不算异常，不打日志
	if err == redis.Nil {
		return nil, false
	}

	// 真正异常才打日志
	if err != nil {
		util.DebugError(fmt.Sprintf("[RedisGet] err: %v, key: %s", err, key))
		return nil, false
	}

	return val, true
}

// RedisExist 判断key是否存在
func RedisExist(ctx context.Context, key string) bool {
	count, err := redisRDB.Exists(ctx, key).Result()
	if err != nil {
		util.DebugError(fmt.Sprintf("[RedisExist] err: %v, key: %s", err, key))
		return false
	}
	return count > 0
}

// RedisSet 标准set（带ctx、带过期）
func RedisSet(ctx context.Context, key string, value any, expiration time.Duration) bool {
	err := redisRDB.Set(ctx, key, value, expiration).Err()
	if err != nil {
		util.DebugError(fmt.Sprintf("[RedisSet] err: %v, key: %s", err, key))
		return false
	}
	return true
}

// RedisChange,只改value
func RedisChange(ctx context.Context, key string, value any) bool {
	err := redisRDB.Set(ctx, key, value, redis.KeepTTL).Err()
	if err != nil {
		util.DebugError(fmt.Sprintf("[RedisChange] err: %v, key: %s", err, key))
		return false
	}
	return true
}

// RedisDel 删除key
func RedisDel(ctx context.Context, key string) bool {
	err := redisRDB.Del(ctx, key).Err()
	if err != nil {
		util.DebugError(fmt.Sprintf("[RedisDel] err: %v, key: %s", err, key))
		return false
	}
	return true
}

// RedisExpire 给key设置过期时间
func RedisExpire(ctx context.Context, key string, expiration time.Duration) bool {

	err := redisRDB.Expire(ctx, key, expiration).Err()
	if err != nil {
		util.DebugError(fmt.Sprintf("[RedisExpire] err: %v, key: %s", err, key))
		return false
	}
	return true
}

func RedisListPush(ctx context.Context, queue string, data string) bool {
	err := redisRDB.LPush(
		ctx,
		queue,
		data,
	).Err()
	if err != nil {
		util.DebugError(fmt.Sprintf("[RedisListPush] err: %v, queue: %s", err, queue))
		return false
	}
	return true
}

func RedisListPop(ctx context.Context, queue string) (string, bool) {
	val, err := redisRDB.RPop(ctx, queue).Result()
	if err == redis.Nil {
		return "", false
	}
	if err != nil {
		util.DebugError(fmt.Sprintf("[RedisListPop] err: %v, queue: %s", err, queue))
		return "", false
	}
	return val, true
}
