package dao

import (
	"newaoe/util"
	"time"

	"xorm.io/xorm"
)

// 代码文件状态的结构体
type CodeAssessmentInfo struct {
	Indices    int       `json:"indices" xorm:"indices pk autoincr indices"` //记录索引
	ID         string    `json:"id" xorm:"id"`                               //用户学号
	UploadTime time.Time `json:"uploadtime" xorm:"uploadtime"`               //提交日期
	Header     string    `json:"header" xorm:"header"`                       //头文件目录(private/id/...)
	Source     string    `json:"source" xorm:"source"`                       //源文件目录(private/id/...)
	Teacher    string    `json:"teacher" xorm:"teacher"`                     //归属教师
	HeaderSize int64     `json:"headersize" xorm:"headersize"`               //头文件大小
	SourceSize int64     `json:"sourcesize" xorm:"sourcesize"`               //源文件大小
}

func (CodeAssessmentInfo) TableName() string {
	return "CodeSubmitAssessment" // 这里返回你想要的表名
}

func CodeAssessmentAdd(session *xorm.Session, data CodeAssessmentInfo) bool {
	_, err := session.Insert(data)
	if err != nil {
		util.Debug("CodeAssessmentAdd:", err)
		return false
	}
	return true
}

func CodeAssessmentGetByID(id string) []CodeAssessmentInfo {
	var ret []CodeAssessmentInfo
	err := DB.Where("id = ?", id).Find(&ret)
	if err != nil {
		util.Debug("CodeAssessmentAddGetByID:", err)
		return nil
	}
	return ret
}
