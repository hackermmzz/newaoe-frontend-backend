package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/redis"
	"newaoe/Src/user/model"
	"newaoe/Src/util"
	"time"
)

func VerifyCodeAdd(keyID string, code string, class int8) error {
	var expire_time time.Duration
	var resend_time time.Duration
	switch class {
	case model.VerifyCode_Class_Regist:
		expire_time = time.Duration(config.Conf.RegistVerifyCode.ExpireTime) * time.Second
		resend_time = time.Duration(config.Conf.RegistVerifyCode.ResendTime) * time.Second
	case model.VerifyCode_Class_PasswordForget:
		expire_time = time.Duration(config.Conf.PasswordForgetVerifyCode.ExpireTime) * time.Second
		resend_time = time.Duration(config.Conf.PasswordForgetVerifyCode.ResendTime) * time.Second
	}

	data := model.VerifyCode{
		KeyID: keyID,
		Class: class,
	}
	tx := redis.RDB.TxPipeline()
	ctx := redis.RDB.Context()
	tx.Set(ctx, data.String(), model.VerifyCodeInRedis{
		Code:       code,
		RetryCount: int8(config.Conf.Other.VerifyCodeRetryCount),
	}.Marshal(), expire_time)
	tx.Set(ctx, data.Tag(), "", resend_time)
	_, err := tx.Exec(ctx)
	if err != nil {
		util.DebugError("VerifyCodeAdd:", err)
		return util.NewError("哈希碰撞!")
	}

	return nil
}

func VerifyCodeExist(keyID string, code string, class int8) (bool, error) {
	data := model.VerifyCode{
		KeyID: keyID,
		Class: class,
	}
	key := data.String()
	byte_data, exist := redis.RedisGet(context.Background(), key)
	if !exist {
		return false, util.NewError("查如此验证码!")
	}
	//
	var info model.VerifyCodeInRedis
	err := json.Unmarshal(byte_data, &info)
	if err != nil {
		return false, err
	}
	//如果验证码不匹配直接返回
	if info.RetryCount <= 0 {
		return false, util.NewError("达到最大可重试次数!")
	}
	//减少1次尝试次数
	redis.RedisChange(context.Background(), key, model.VerifyCodeInRedis{
		Code:       info.Code,
		RetryCount: info.RetryCount - 1,
	}.Marshal())
	//验证码错误
	if code != info.Code {
		return false, util.NewError(fmt.Sprintf("验证码错误,你还有%d尝试", info.RetryCount-1))
	}
	//
	return true, nil
}

func VerifyCodeCanSendTTL(keyID string, class int8) int64 {
	data := model.VerifyCode{
		KeyID: keyID,
		Class: class,
	}
	ttl, err := redis.RDB.TTL(redis.RDB.Context(), data.Tag()).Result()
	if err != nil {
		util.DebugError("VerifyCodeCanSendTTL:", err)
		return int64(1e9)
	}
	return int64(ttl.Seconds())
}

func VerifyCodeRemove(keyID string, class int8) bool {
	data := model.VerifyCode{
		KeyID: keyID,
		Class: class,
	}
	ok0 := redis.RedisDel(context.Background(), data.String())
	//ok1 := redis.RedisDel(context.Background(), data.Tag())//冷却不删
	return ok0
}

// 注册
func RegistVerifyCodeAdd(keyID string, code string) error {
	return VerifyCodeAdd(keyID, code, model.VerifyCode_Class_Regist)
}

func RegistVerifyCodeExist(keyID string, code string) (bool, error) {
	return VerifyCodeExist(keyID, code, model.VerifyCode_Class_Regist)
}

func RegistVerifyCodeCanSendTTL(keyID string) int64 {
	return VerifyCodeCanSendTTL(keyID, model.VerifyCode_Class_Regist)
}

func RegistVerifyCodeRemove(keyID string) bool {
	return VerifyCodeRemove(keyID, model.VerifyCode_Class_Regist)
}

// 忘记密码重置验证码
func PasswordForgetVerifyCodeAdd(keyID string, code string) error {
	return VerifyCodeAdd(keyID, code, model.VerifyCode_Class_PasswordForget)
}

func PasswordForgetVerifyCodeExist(keyID string, code string) (bool, error) {
	return VerifyCodeExist(keyID, code, model.VerifyCode_Class_PasswordForget)
}

func PasswordForgetVerifyCodeCanSendTTL(keyID string) int64 {
	return VerifyCodeCanSendTTL(keyID, model.VerifyCode_Class_PasswordForget)
}
func PasswordForgetVerifyCodeRemove(keyID string) bool {
	return VerifyCodeRemove(keyID, model.VerifyCode_Class_PasswordForget)
}
