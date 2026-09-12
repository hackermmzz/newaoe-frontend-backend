package data

import (
	"fmt"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"
	"time"
)

// 定义verifyCode类型
const (
	VerifyCode_Class_Regist         = 0 //注册验证码
	VerifyCode_Class_PasswordForget = 1 //登录验证码
)

// Define the VerifyCode struct
type VerifyCode struct {
	Code  string
	Id    string
	Class int8
}

// 用于检查验证码是否存在
func (v VerifyCode) String() string {
	return fmt.Sprintf("VerifyCode:%v:ID:%v:Class:%v", v.Code, v.Id, v.Class)
}

// 用于检查验证码是否可以重发
func (v VerifyCode) Tag() string {
	return fmt.Sprintf("VerifyCode::ID:%v:Class:%v", v.Id, v.Class)
}

func VerifyCodeAdd(id string, code string, class int8) bool {
	var expire_time time.Duration
	var resend_time time.Duration
	switch class {
	case VerifyCode_Class_Regist:
		expire_time = time.Duration(config.Conf.RegistVerifyCode.ExpireTime) * time.Second
		resend_time = time.Duration(config.Conf.RegistVerifyCode.ResendTime) * time.Second
	case VerifyCode_Class_PasswordForget:
		expire_time = time.Duration(config.Conf.PasswordForgetVerifyCode.ExpireTime) * time.Second
		resend_time = time.Duration(config.Conf.PasswordForgetVerifyCode.ResendTime) * time.Second
	}

	data := VerifyCode{
		Id:    id,
		Code:  code,
		Class: class,
	}
	tx := dao.RDB.TxPipeline()
	ctx := dao.RDB.Context()
	tx.Set(ctx, data.String(), "", expire_time)
	tx.Set(ctx, data.Tag(), "", resend_time)
	_, err := tx.Exec(ctx)
	if err != nil {
		util.Debug("VerifyCodeAdd:", err)
		return false
	}

	return true
}

func VerifyCodeExist(id string, code string, class int8) bool {
	data := VerifyCode{
		Id:    id,
		Code:  code,
		Class: class,
	}
	exist, err := dao.RDB.Exists(dao.RDB.Context(), data.String()).Result()
	if err != nil {
		util.Debug("VerifyCodeExist:", err)
		return false
	}
	return exist == 1
}

func VerifyCodeCanSend(id string, class int8) bool {
	data := VerifyCode{
		Id:    id,
		Class: class,
	}
	exist, err := dao.RDB.Exists(dao.RDB.Context(), data.Tag()).Result()
	if err != nil {
		util.Debug("VerifyCodeCanSend:", err)
		return false
	}
	return exist == 1
}

func VerifyCodeCanSendTTL(id string, class int8) int64 {
	data := VerifyCode{
		Id:    id,
		Class: class,
	}
	ttl, err := dao.RDB.TTL(dao.RDB.Context(), data.Tag()).Result()
	if err != nil {
		util.Debug("VerifyCodeCanSendTTL:", err)
		return int64(1e9)
	}
	return int64(ttl.Seconds())
}

// 注册
func RegistVerifyCodeAdd(id string, code string) bool {
	return VerifyCodeAdd(id, code, VerifyCode_Class_Regist)
}

func RegistVerifyCodeExist(id string, code string) bool {
	return VerifyCodeExist(id, code, VerifyCode_Class_Regist)
}

func RegistVerifyCodeCanSend(id string) bool {
	return VerifyCodeCanSend(id, VerifyCode_Class_Regist)
}
func RegistVerifyCodeCanSendTTL(id string) int64 {
	return VerifyCodeCanSendTTL(id, VerifyCode_Class_Regist)
}

// 忘记密码重置验证码
func PasswordForgetVerifyCodeAdd(id string, code string) bool {
	return VerifyCodeAdd(id, code, VerifyCode_Class_PasswordForget)
}

func PasswordForgetVerifyCodeExist(id string, code string) bool {
	return VerifyCodeExist(id, code, VerifyCode_Class_PasswordForget)
}

func PasswordForgetVerifyCodeCanSend(id string) bool {
	return VerifyCodeCanSend(id, VerifyCode_Class_PasswordForget)
}
func PasswordForgetVerifyCodeCanSendTTL(id string) int64 {
	return VerifyCodeCanSendTTL(id, VerifyCode_Class_PasswordForget)
}
