package dao

import "time"

type EmailWaitInfo struct {
	Indices      int       `json:"indices" xorm:"pk 'indices'"`      //对应的索引
	LastSendTime time.Time `json:"lastsendtime" xorm:"lastsendtime"` //上一次尝试发送的时间
}
