package model

import "time"

// VIP等级
var (
	VIP_TOURIST = 1000
	VIP_NONE    = 2000
	VIP_SUPER   = 3000
)

type Student struct {
	ID         string    `json:"id" xorm:"id"`
	Password   string    `json:"password" xorm:"password"`
	Email      string    `json:"email" xorm:"email"`
	Avatar     string    `json:"avatar" xorm:"avatar"`
	RegistDate time.Time `json:"registDate" xorm:"registDate"`
	Vip        int       `json:"vip" xorm:"vip"`
}

func (Student) TableName() string {
	return "Student" // 这里返回你想要的表名
}
