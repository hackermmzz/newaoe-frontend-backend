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

// 全局变量
const KeepTTL = redis.KeepTTL

// 全局 Redis 客户端（单例，线程安全）
var redisRDB *redis.Client

// ConnectRedis 生产级 Redis 初始化（连接池 + 重试 + 健康检查）
func ConnectRedis() error {
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
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("Redis 连接最终失败: %w", err)
	}

	util.DebugSuccess("Redis连接池初始化成功")
	return nil
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
func RedisGet(ctx context.Context, key string) ([]byte, bool, error) {
	val, err := redisRDB.Get(ctx, key).Bytes()

	// key不存在不算异常，不打日志
	if err == redis.Nil {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return val, true, nil
}

// RedisExist 判断key是否存在
func RedisExist(ctx context.Context, key string) (bool, error) {
	count, err := redisRDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// RedisSet 标准set（带ctx、带过期）
func RedisSet(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	err := redisRDB.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

// 自增
func RedisIncrease(ctx context.Context, key string) (int, error) {
	n, err := redisRDB.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

// 自减
func RedisDecrease(ctx context.Context, key string) (int, error) {
	n, err := redisRDB.Decr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

// RedisChange,只改value
func RedisChange(ctx context.Context, key string, value any) (bool, error) {
	err := redisRDB.Set(ctx, key, value, redis.KeepTTL).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

// RedisDel 删除key
func RedisDel(ctx context.Context, key string) (bool, error) {
	err := redisRDB.Del(ctx, key).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

// RedisExpire 给key设置过期时间
func RedisExpire(ctx context.Context, key string, expiration time.Duration) (bool, error) {

	err := redisRDB.Expire(ctx, key, expiration).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func RedisListPush(ctx context.Context, queue string, data string) (bool, error) {
	err := redisRDB.LPush(
		ctx,
		queue,
		data,
	).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func RedisListPop(ctx context.Context, queue string) (string, bool, error) {
	val, err := redisRDB.RPop(ctx, queue).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}
