package model

import (
	"time"
)

// 这个表负责维护哪些正在跑，但是可能会出问题的运行记录（保证失败了可以重新推入队列）
type CodeRunningInfo struct {
	Indices    int       `json:"indices" xorm:"indices pk "`
	SubmitTime time.Time `json:"submittime" xorm:"submittime"` //提交日期
	RunType    int       `json:"runtype" xorm:"runtype"`       //运行类型
	Status     int       `json:"status" xorm:"status"`         //运行阶段
	RetryCount uint8     `json:"retrycount" xorm:"retrycount"` //重试次数
}

func (CodeRunningInfo) TableName() string {
	return "CodeRunning" // 这里返回你想要的表名
}
