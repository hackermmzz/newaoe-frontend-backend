package dao

import (
	"newaoe/util"
	"time"

	"xorm.io/xorm"
)

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

func EmailInfoInsert(session *xorm.Session, info EmailInfo) int {
	// Insert 返回影响行数 + error
	_, err := session.Insert(&info)
	if err != nil {
		util.DebugError("EmailInfoInsert:", err)
		return 0
	}
	return info.Indices
}

// 前一个为是否成功，后一个为是否有改变
func EmailInfoUpdateSendStatus(session *xorm.Session, indices int, send bool) (bool, int64) {
	affected, err := session.ID(indices).
		Cols("send"). // 明确指定只更新send字段
		Update(&EmailInfo{
			Send: send,
		})
	if err != nil {
		util.DebugError("EmailUpdateSendStatus:", err)
		return false, 0
	}
	return true, affected
}

func EmailInfoGetForUpdate(session *xorm.Session, indices int) *EmailInfo {
	var ret EmailInfo
	has, err := session.ForUpdate().Where("indices = ?", indices).Get(&ret)
	if !has {
		return nil
	}
	if err != nil {
		util.DebugError("EmailInfoGet:", err)
		return nil
	}
	return &ret
}
