package dao

import (
	"newaoe/util"
	"time"

	"xorm.io/xorm"
)

type FeedbackInfo struct {
	Indices    int       `json:"indices" xorm:"pk autoincr 'indices'"`
	SubmitTime time.Time `json:"submittime" xorm:"notnull 'submittime'"`
	BaseFolder string    `json:"basefolder" xorm:"basefolder"`
}

func (c FeedbackInfo) TableName() string {
	return "Feedback"
}

func FeedbackGetByIndices(session *xorm.Session, indices int) *FeedbackInfo {
	var ret FeedbackInfo
	has, err := session.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		util.DebugError("FeedbackGetByIndices:", err)
		return nil
	}
	return &ret
}

func FeedbackInsert(session *xorm.Session, info FeedbackInfo) int {
	// Insert 返回影响行数 + error
	_, err := session.Insert(&info)
	if err != nil {
		util.DebugError("FeedbackInsert:", err)
		return 0
	}
	return info.Indices
}
