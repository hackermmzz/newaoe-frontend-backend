package dao

import (
	"newaoe/util"
	"time"

	"xorm.io/xorm"
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

func CodeRunAdd(session *xorm.Session, c CodeRunInfo) int {
	_, err := session.Insert(&c)
	if err != nil {
		util.Debug("CodeRunAdd:", err)
		return 0
	}
	return c.Indices
}

func CodeRunUpdateStatusAsync(session *xorm.Session, indices int, status string) bool {
	var info CodeRunInfo
	info.Status = status
	_, err := session.Where("indices=?", indices).
		Update(info)
	if err != nil {
		util.Debug("CodeRunUpdateStatusAsync:", err)
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
		util.Debug("CodeRunUpdateStatusSync:", err)
		return false
	}
	return true
}

func CodeRunGetByIndices(indices int) *CodeRunInfo {
	var ret CodeRunInfo
	has, err := DB.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		util.Debug("CodeRunGetByIndices:", err)
		return nil
	}
	return &ret
}

func CodeRunGetById(id string) []CodeRunInfo {
	var ret []CodeRunInfo
	err := DB.Where("id = ?", id).Find(&ret)
	if err != nil {
		util.Debug("CodeRunGetById:", err)
		return nil
	}
	return ret
}

func CodeRunCountById(id string) int64 {
	cnt, err := DB.Where("id = ?", id).Count(&CodeRunInfo{})
	if err != nil {
		util.Debug("CodeRunCountById:", err)
		return 0
	}
	return cnt
}

func CodeRunGetRangeById(id string, beg int, end int) []CodeRunInfo {
	var ret []CodeRunInfo
	err := DB.Where("id = ?", id).
		OrderBy("submittime desc").
		Limit(end-beg+1, int(beg)).
		Find(&ret)
	if err != nil {
		util.Debug("CodeRunGetRangeById:", err)
		return nil
	}
	return ret
}
