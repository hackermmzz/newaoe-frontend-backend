package model

import (
	"encoding/json"
	"time"
)

const (
	AssessmenrRankStatus_Error       = -2
	AssessmenrRankStatus_CompileFail = -1
	AssessmenrRankStatus_Crash       = 0
	AssessmenrRankStatus_Fail        = 1
	AssessmenrRankStatus_Win         = 2
)

type AssessmentRankMsgInfo struct {
	Status CodeRunStatusInfo `json:"status"`
}

func (c AssessmentRankMsgInfo) Marshal() string {
	d, _ := json.Marshal(c)
	return string(d)
}

func (c *AssessmentRankMsgInfo) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &c)
	return err
}

type AssessmentRankInfo struct {
	ID         string    `json:"id" xorm:"id pk"`              //用户学号
	Status     int       `json:"status" xorm:"status"`         //胜利或者失败或者崩溃
	SubmitTime time.Time `json:"submittime" xorm:"submittime"` //提交日期
	Score      int       `json:"score" xorm:"score"`           //分数
	Frame      int       `json:"frame" xorm:"frame"`           //运行时间
	Msg        string    `json:"msg" xorm:"msg"`               //运行结果(AssessmentRankMsgInfo)
}

func (r AssessmentRankInfo) TableName() string {
	return "AssessmentRank"
}
