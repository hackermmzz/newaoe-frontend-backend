package model

import "time"

// VIP等级
var (
	VIP_NONE  = 0 //没有特权
	VIP_SUPER = 1 //特权用户
)

type Student struct {
	Id         string    `json:"id" xorm:"id"`
	Password   string    `json:"password" xorm:"password"`
	Email      string    `json:"email" xorm:"email"`
	Avatar     string    `json:"avatar" xorm:"avatar"`
	RegistDate time.Time `json:"registDate" xorm:"registDate"`
	Vip        int       `json:"vip" xorm:"vip"`
}

func (Student) TableName() string {
	return "Student" // 这里返回你想要的表名
}
