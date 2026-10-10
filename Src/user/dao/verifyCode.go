package dao

import (
	"context"
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/redis"
	"newaoe/Src/user/model"
	"newaoe/Src/util"
	"time"
)

var verifyCodeLua = `
local value = redis.call("GET", KEYS[1])

if not value then
    return {-1, 0}
end

local info = cjson.decode(value)

if info.retrycount <= 0 then
    return {-2, 0}
end

if ARGV[1] == info.code then
    return {1, info.retrycount}
end

info.retrycount = info.retrycount - 1

redis.call(
    "SET",
    KEYS[1],
    cjson.encode(info),
    "KEEPTTL"
)

return {0, info.retrycount}
`

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

	ctx := context.Background()
	err := redis.RedisTx(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, data.String(), model.VerifyCodeInRedis{
			Code:       code,
			RetryCount: int8(config.Conf.Other.VerifyCodeRetryCount),
		}.Marshal(), expire_time)
		pipe.Set(ctx, data.Tag(), "", resend_time)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func VerifyCodeExist(keyID string, code string, class int8) (bool, error) {
	data := model.VerifyCode{
		KeyID: keyID,
		Class: class,
	}
	//lua执行
	resultCMD := redis.RedisLua(
		context.Background(),
		verifyCodeLua,
		[]string{data.String()},
		code,
	)
	//获取结果
	result, err := resultCMD.Int64Slice()
	if err != nil {
		return false, err
	}
	switch result[0] {
	case -1:
		return false, util.NewError("无此验证码!")
	case -2:
		return false, util.NewError("达到最大可重试次数!")
	case 0:
		return false, util.NewError(fmt.Sprintf("验证码错误,你还有%d次尝试", result[1]))
	case 1:
		return true, nil
	}
	return false, util.NewError("验证码校验异常")
}

func VerifyCodeCanSendTTL(keyID string, class int8) (int64, error) {
	data := model.VerifyCode{
		KeyID: keyID,
		Class: class,
	}
	ttl, err := redis.RedisTTL(context.Background(), data.Tag())
	if err != nil {
		return int64(1e9), err
	}
	return int64(ttl.Seconds()), nil
}

func VerifyCodeRemove(keyID string, class int8) (bool, error) {
	data := model.VerifyCode{
		KeyID: keyID,
		Class: class,
	}
	ok0, err := redis.RedisDel(context.Background(), data.String())
	if err != nil {
		return false, err
	}
	//ok1 := redis.RedisDel(context.Background(), data.Tag())//冷却不删
	return ok0, nil
}

// 注册
func RegistVerifyCodeAdd(keyID string, code string) error {
	return VerifyCodeAdd(keyID, code, model.VerifyCode_Class_Regist)
}

func RegistVerifyCodeExist(keyID string, code string) (bool, error) {
	return VerifyCodeExist(keyID, code, model.VerifyCode_Class_Regist)
}

func RegistVerifyCodeCanSendTTL(keyID string) (int64, error) {
	return VerifyCodeCanSendTTL(keyID, model.VerifyCode_Class_Regist)
}

func RegistVerifyCodeRemove(keyID string) (bool, error) {
	return VerifyCodeRemove(keyID, model.VerifyCode_Class_Regist)
}

// 忘记密码重置验证码
func PasswordForgetVerifyCodeAdd(keyID string, code string) error {
	return VerifyCodeAdd(keyID, code, model.VerifyCode_Class_PasswordForget)
}

func PasswordForgetVerifyCodeExist(keyID string, code string) (bool, error) {
	return VerifyCodeExist(keyID, code, model.VerifyCode_Class_PasswordForget)
}

func PasswordForgetVerifyCodeCanSendTTL(keyID string) (int64, error) {
	return VerifyCodeCanSendTTL(keyID, model.VerifyCode_Class_PasswordForget)
}
func PasswordForgetVerifyCodeRemove(keyID string) (bool, error) {
	return VerifyCodeRemove(keyID, model.VerifyCode_Class_PasswordForget)
}
