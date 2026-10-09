package model

import "time"

// 运行类型（Release运行或者重度Debug运行）
var (
	CodeRunType_ReleaseRun = 0 //release版本运行
	CodeRunType_DebugRun   = 1 //重度调试模式运行
)

// 上传的代码类别
var (
	Code_CommonSubmit     = 1 //普通提交
	Code_AssessmentSubmit = 2 //考核提交
	Code_ReRunSubmit      = 3 //重新运行
)

// 这里我们采用分表的方式来进行
type CodeRunInfo struct {
	Indices     int       `json:"indices" xorm:"indices pk autoincr indices"` //记录索引
	ID          string    `json:"id" xorm:"id"`                               //用户学号
	SubmitTime  time.Time `json:"submittime" xorm:"submittime"`               //提交日期
	Header      string    `json:"header" xorm:"header"`                       //头文件目录(private/id/...)
	Source      string    `json:"source" xorm:"source"`                       //源文件目录(private/id/...)
	Description string    `json:"description" xorm:"description"`             //描述
	Class       int       `json:"class" xorm:"class"`                         //代码类型
	Status      string    `json:"status" xorm:"status"`                       //运行状态
	Version     int64     `json:"version" xorm:"version"`                     //版本(防止旧数据覆盖)
}

func (c CodeRunInfo) TableName() string {
	return "CodeRun"
}
