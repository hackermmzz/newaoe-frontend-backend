package model

import (
	"fmt"
)

// 定义verifyCode类型
const (
	VerifyCode_Class_Regist         = 0 //注册验证码
	VerifyCode_Class_PasswordForget = 1 //登录验证码
)

// Define the VerifyCode struct
type VerifyCode struct {
	Code  string `json:"code"`
	KeyID string `json:"keyid"`
	Class int8   `json:"class"`
}

// 用于检查验证码是否存在
func (v VerifyCode) String() string {
	return fmt.Sprintf("VerifyCode:%v:ID:%v:Class:%v", v.Code, v.KeyID, v.Class)
}

// 用于检查验证码是否可以重发
func (v VerifyCode) Tag() string {
	return fmt.Sprintf("VerifyCode::ID:%v:Class:%v", v.KeyID, v.Class)
}
