package model

import "time"

type RankInfo struct {
	ID         string    `json:"id" xorm:"id pk"`              //用户学号
	Win        bool      `json:"win" xorm:"win"`               //胜利或者失败(娄老师说只能是成功的才能进榜)
	SubmitTime time.Time `json:"submittime" xorm:"submittime"` //提交日期
	Score      int       `json:"score" xorm:"score"`           //分数
	Frame      int       `json:"frame" xorm:"frame"`           //运行时间
	Msg        string    `json:"msg" xorm:"msg"`               //运行结果
}

func (r RankInfo) TableName() string {
	return "Rank"
}
