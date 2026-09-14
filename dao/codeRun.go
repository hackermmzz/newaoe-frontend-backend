package dao

import (
	"encoding/json"
	"newaoe/util"
	"time"

	"xorm.io/xorm"
)

// 代码当前的运行状态
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

func CodeRunAdd(session *xorm.Session, c CodeRunInfo) int {
	_, err := session.Insert(&c)
	if err != nil {
		util.DebugError("CodeRunAdd:", err)
		return 0
	}
	return c.Indices
}

func CodeRunUpdate(session *xorm.Session, newInfo CodeRunInfo) bool {
	affected, err := session.
		ID(newInfo.Indices).
		AllCols().
		Update(&newInfo)
	if err != nil {
		util.DebugError("CodeRunUpdate:", err)
		return false
	}
	return affected > 0
}

func CodeRunUpdateStatusAsync(session *xorm.Session, indices int, status string) bool {
	var info CodeRunInfo
	info.Status = status
	_, err := session.Where("indices=?", indices).
		Update(info)
	if err != nil {
		util.DebugError("CodeRunUpdateStatusAsync:", err)
		return false
	}
	return true
}

func CodeRunUpdateStatusSync(indices int, status string) bool {
	var info CodeRunInfo
	info.Status = status
	_, err := DB.Where("indices=?", indices).
		Update(info)
	if err != nil {
		util.DebugError("CodeRunUpdateStatusSync:", err)
		return false
	}
	return true
}

func CodeRunGetByIndices(indices int) *CodeRunInfo {
	var ret CodeRunInfo
	has, err := DB.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		util.DebugError("CodeRunGetByIndices:", err)
		return nil
	}
	return &ret
}

func CodeRunBatchGetByIndices(indices []int) []CodeRunInfo {
	var ret []CodeRunInfo
	if len(indices) == 0 {
		return ret
	}
	err := DB.
		Where("indices IN (?)", indices).
		Find(&ret)
	if err != nil {
		util.DebugError("CodeRunBatchGetByIndices:", err)
		return nil
	}
	return ret
}

func CodeRunGetById(id string) []CodeRunInfo {
	var ret []CodeRunInfo
	err := DB.Where("id = ?", id).Find(&ret)
	if err != nil {
		util.DebugError("CodeRunGetById:", err)
		return nil
	}
	return ret
}

func CodeRunCountById(id string) int64 {
	cnt, err := DB.Where("id = ?", id).Count(&CodeRunInfo{})
	if err != nil {
		util.DebugError("CodeRunCountById:", err)
		return 0
	}
	return cnt
}

// 不包含end
func CodeRunGetRangeById(id string, beg int, end int) []CodeRunInfo {
	if beg < 0 || beg >= end {
		return nil
	}
	var ret []CodeRunInfo
	err := DB.Where("id = ?", id).
		OrderBy("submittime desc").
		Limit(end-beg, int(beg)).
		Find(&ret)
	if err != nil {
		util.DebugError("CodeRunGetRangeById:", err)
		return nil
	}
	return ret
}
