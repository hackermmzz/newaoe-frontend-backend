package model

import (
	"encoding/json"
	"time"
)

// 代码当前的运行状态(不要随便改动顺序)
var (
	Code_Status_Error           = 0 //服务器异常
	Code_Status_Wait            = 1 //在等待队列里面
	Code_Status_Compile         = 2 //编译中
	Code_Status_Compile_Success = 3 //编译成功
	Code_Status_Compile_Fail    = 4 //编译错误
	Code_Status_Running         = 5 //正在运行出结果
	Code_Status_Success         = 6 //运行胜利
	Code_Status_Fail            = 7 //游戏失败
	Code_Status_Crash           = 8 //游戏崩溃

)

type CodeRunStatusInfo struct {
	Status int    `json:"status"`
	Food   int    `json:"food"`
	Wood   int    `json:"wood"`
	Gold   int    `json:"gold"`
	Stone  int    `json:"stone"`
	Frame  int    `json:"frame"`
	Win    bool   `json:"win"`
	Score  int    `json:"score"`
	Data   string `json:"data"`
}

// 这个表负责维护哪些正在跑，但是可能会出问题的运行记录（保证失败了可以重新推入队列）
type CodeRunningInfo struct {
	Indices    int       `json:"indices" xorm:"indices pk "`
	SubmitTime time.Time `json:"submittime" xorm:"submittime"` //提交日期
	RunType    int       `json:"runtype" xorm:"runtype"`       //运行类型
}

func (CodeRunningInfo) TableName() string {
	return "CodeRunning" // 这里返回你想要的表名
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
	case Code_Status_Error, Code_Status_Compile_Fail,
		Code_Status_Success, Code_Status_Fail, Code_Status_Crash:
		return true
	}
	return false
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
