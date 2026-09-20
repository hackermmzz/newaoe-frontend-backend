package dao

import (
	"newaoe/Src/config"
	"newaoe/Src/redis"
	"newaoe/Src/user/model"
	"newaoe/Src/util"
	"time"
)

func VerifyCodeAdd(keyID string, code string, class int8) bool {
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
		Code:  code,
		Class: class,
	}
	tx := redis.RDB.TxPipeline()
	ctx := redis.RDB.Context()
	tx.Set(ctx, data.String(), "", expire_time)
	tx.Set(ctx, data.Tag(), "", resend_time)
	_, err := tx.Exec(ctx)
	if err != nil {
		util.DebugError("VerifyCodeAdd:", err)
		return false
	}

	return true
}

func VerifyCodeExist(keyID string, code string, class int8) bool {
	data := model.VerifyCode{
		KeyID: keyID,
		Code:  code,
		Class: class,
	}
	exist, err := redis.RDB.Exists(redis.RDB.Context(), data.String()).Result()
	if err != nil {
		util.DebugError("VerifyCodeExist:", err)
		return false
	}
	return exist == 1
}

func VerifyCodeCanSend(keyID string, class int8) bool {
	data := model.VerifyCode{
		KeyID: keyID,
		Class: class,
	}
	exist, err := redis.RDB.Exists(redis.RDB.Context(), data.Tag()).Result()
	if err != nil {
		util.DebugError("VerifyCodeCanSend:", err)
		return false
	}
	return exist == 1
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

// 注册
func RegistVerifyCodeAdd(keyID string, code string) bool {
	return VerifyCodeAdd(keyID, code, model.VerifyCode_Class_Regist)
}

func RegistVerifyCodeExist(keyID string, code string) bool {
	return VerifyCodeExist(keyID, code, model.VerifyCode_Class_Regist)
}

func RegistVerifyCodeCanSend(keyID string) bool {
	return VerifyCodeCanSend(keyID, model.VerifyCode_Class_Regist)
}
func RegistVerifyCodeCanSendTTL(keyID string) int64 {
	return VerifyCodeCanSendTTL(keyID, model.VerifyCode_Class_Regist)
}

// 忘记密码重置验证码
func PasswordForgetVerifyCodeAdd(keyID string, code string) bool {
	return VerifyCodeAdd(keyID, code, model.VerifyCode_Class_PasswordForget)
}

func PasswordForgetVerifyCodeExist(keyID string, code string) bool {
	return VerifyCodeExist(keyID, code, model.VerifyCode_Class_PasswordForget)
}

func PasswordForgetVerifyCodeCanSend(keyID string) bool {
	return VerifyCodeCanSend(keyID, model.VerifyCode_Class_PasswordForget)
}
func PasswordForgetVerifyCodeCanSendTTL(keyID string) int64 {
	return VerifyCodeCanSendTTL(keyID, model.VerifyCode_Class_PasswordForget)
}
