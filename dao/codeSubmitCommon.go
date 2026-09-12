package dao

import (
	"newaoe/util"
	"time"

	"xorm.io/xorm"
)

// 上传的代码类别
var (
	Code_Common                     = 1 //普通提交
	Code_AssessmentSubmissionUpload = 2 //考核提交
)

// 代码文件状态的结构体
type CodeCommonInfo struct {
	Indices     int       `json:"indices" xorm:"indices pk autoincr indices"` //记录索引
	ID          string    `json:"id" xorm:"id"`                               //用户学号
	UploadTime  time.Time `json:"uploadtime" xorm:"uploadtime"`               //提交日期
	Header      string    `json:"header" xorm:"header"`                       //头文件目录(private/id/...)
	Source      string    `json:"source" xorm:"source"`                       //源文件目录(private/id/...)
	Description string    `json:"description" xorm:"description"`             //描述
	HeaderSize  int64     `json:"headersize" xorm:"headersize"`               //头文件大小
	SourceSize  int64     `json:"sourcesize" xorm:"sourcesize"`               //源文件大小
}

func (CodeCommonInfo) TableName() string {
	return "CodeSubmitCommon" // 这里返回你想要的表名
}

func CodeCommonInfoAdd(session *xorm.Session, data CodeCommonInfo) bool {
	_, err := session.Insert(data)
	if err != nil {
		util.Debug("CodeCommonInfoAdd:", err)
		return false
	}
	return true
}

func CodeCommonGetByIndices(indices int) *CodeCommonInfo {
	var ret CodeCommonInfo
	has, err := DB.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		util.Debug("CodeCommonGetByIndices:", err)
		return nil
	}
	return &ret
}
func CodeCommonGetByID(id string) []CodeCommonInfo {
	var ret []CodeCommonInfo
	err := DB.Where("id = ?", id).Find(&ret)
	if err != nil {
		util.Debug("CodeCommonGetByID:", err)
		return nil
	}
	return ret
}
