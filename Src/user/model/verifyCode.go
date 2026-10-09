package model

import (
	"encoding/json"
	"fmt"
)

// 定义verifyCode类型
const (
	VerifyCode_Class_Regist         = 0 //注册验证码
	VerifyCode_Class_PasswordForget = 1 //登录验证码
)

// Define the VerifyCode struct
type VerifyCode struct {
	KeyID string `json:"keyid"`
	Class int8   `json:"class"`
}

// 存在redis里面的数据
type VerifyCodeInRedis struct {
	Code       string `json:"code"`
	RetryCount int8   `json:"retrycount"`
}

// 验证码的Key
func (v VerifyCode) String() string {
	return fmt.Sprintf("VerifyCode:KeyID:%v:Class:%v", v.KeyID, v.Class)
}

// 判断重发的Key
func (v VerifyCode) Tag() string {
	return fmt.Sprintf("VerifyCodeRsend:KeyID:%v:Class:%v", v.KeyID, v.Class)
}

// 序列化
func (v VerifyCodeInRedis) Marshal() []byte {
	dt, _ := json.Marshal(v)
	return dt
}
