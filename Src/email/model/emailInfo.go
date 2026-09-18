package model

import "time"

var (
	Email_RegistCode               = 0 //注册码
	Email_PasswordForgetVerifyCode = 1 //忘记密码验证码
	Email_Feedback                 = 2 //反馈
)

type EmailInfo struct {
	Indices    int       `json:"indices" xorm:"pk autoincr 'indices'"`
	Receiver   string    `json:"receiver" xorm:"notnull 'receiver'"`
	CreateTime time.Time `json:"createtime" xorm:"notnull 'createtime'"`
	Class      int       `json:"class" xorm:"notnull 'class'"`
	Data       string    `json:"data" xorm:"text notnull 'data'"`
	Send       bool      `json:"send" xorm:"notnull 'send'"`
}

func (r EmailInfo) TableName() string {
	return "EmailInfo"
}
