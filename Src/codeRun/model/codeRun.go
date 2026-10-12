package model

import (
	"encoding/json"
	"time"
)

/********************************运行状态表*********************************/
// 代码当前的运行状态(不要随便改动顺序,这个是经过排序的)
var (
	Code_Status_Unknown         = 0    //未知状态
	Code_Status_Error           = 1000 //服务器异常
	Code_Status_Wait            = 2000 //在等待队列里面
	Code_Status_Compile         = 3000 //编译中
	Code_Status_Compile_Fail    = 4000 //编译错误
	Code_Status_Compile_Success = 5000 //编译成功
	Code_Status_Running         = 6000 //正在运行出结果
	Code_Status_Crash           = 7000 //游戏崩溃
	Code_Status_Fail            = 8000 //游戏失败
	Code_Status_Success         = 9000 //运行胜利
)

type CodeRunStatusInfo struct {
	Status   int    `json:"status"`
	Food     int    `json:"food"`
	Wood     int    `json:"wood"`
	Gold     int    `json:"gold"`
	Stone    int    `json:"stone"`
	Frame    int    `json:"frame"`
	Win      bool   `json:"win"`
	Score    int    `json:"score"`
	Data     string `json:"data"`
	TimeCost int64  `json:"time_cost"` //运行耗时s
}

func NewCodeRunStatusInfo() CodeRunStatusInfo {
	return CodeRunStatusInfo{
		Status: Code_Status_Wait,
	}
}
func (c CodeRunStatusInfo) Marshal() string {
	d, _ := json.Marshal(c)
	return string(d)
}

func (c *CodeRunStatusInfo) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &c)
	return err
}
func (c CodeRunStatusInfo) IsFinish() bool {
	switch c.Status {
	case Code_Status_Compile_Fail, Code_Status_Success, Code_Status_Fail, Code_Status_Crash:
		return true
	}
	return false
}

// 判断当前的status是否比oldStatus好
func (p CodeRunStatusInfo) IsBetter(oldStatus CodeRunStatusInfo) bool {
	if p.Status != oldStatus.Status {
		return p.Status > oldStatus.Status
	}
	//状态相同，比较的得分
	if p.Win {
		//两个都是胜利，先比较帧率
		if p.Frame < oldStatus.Frame {
			return true
		}
		//帧率相同，比较得分
		return p.Score > oldStatus.Score
	} else {
		//都是失败，先比较得分
		if p.Score > oldStatus.Score {
			return true
		}
		//再比较运行帧率
		return p.Frame > oldStatus.Frame
	}
}

func ProcessDataMessageByStatus(coderunStatus CodeRunStatusInfo) CodeRunStatusInfo {
	var ret CodeRunStatusInfo
	ret = coderunStatus
	data := coderunStatus.Data
	switch coderunStatus.Status {
	case Code_Status_Wait:
		ret.Data = "排队中..."
		ret.Status = Code_Status_Wait
	case Code_Status_Compile:
		ret.Data = "编译中..."
		ret.Status = Code_Status_Compile
	case Code_Status_Compile_Fail:
		ret.Data = "编译错误: " + data
		ret.Status = Code_Status_Compile_Fail
	case Code_Status_Compile_Success:
		ret.Data = "编译成功"
		ret.Status = Code_Status_Compile_Success
	case Code_Status_Crash:
		ret.Data = "运行崩溃: " + data
		ret.Status = Code_Status_Crash
	case Code_Status_Fail:
		ret.Data = "游戏失败: " + data
		ret.Status = Code_Status_Fail
	case Code_Status_Success:
		ret.Data = "游戏胜利: " + data
		ret.Status = Code_Status_Success
	case Code_Status_Running:
		ret.Data = "正在奋战: " + data
		ret.Status = Code_Status_Running
	default:
		ret.Data = "服务器异常"
		ret.Status = Code_Status_Error
	}
	return ret
}

/********************************运行表*********************************/
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
}

func (c CodeRunInfo) TableName() string {
	return "CodeRun"
}
