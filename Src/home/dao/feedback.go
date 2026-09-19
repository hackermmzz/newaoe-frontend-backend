package dao

import (
	database "newaoe/Src/databse"
	"newaoe/Src/home/model"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func FeedbackGetByIndices(session *xorm.Session, indices int) *model.FeedbackInfo {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var ret model.FeedbackInfo
	has, err := session.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		util.DebugError("FeedbackGetByIndices:", err)
		return nil
	}
	return &ret
}

func FeedbackInsert(session *xorm.Session, info model.FeedbackInfo) int {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	// Insert 返回影响行数 + error
	_, err := session.Insert(&info)
	if err != nil {
		util.DebugError("FeedbackInsert:", err)
		return 0
	}
	return info.Indices
}
