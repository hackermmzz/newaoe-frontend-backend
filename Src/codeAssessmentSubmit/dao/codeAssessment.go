package dao

import (
	"newaoe/Src/codeAssessmentSubmit/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func CodeAssessmentAdd(session *xorm.Session, data model.CodeAssessmentInfo) int64 {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	_, err := session.Insert(&data)
	if err != nil {
		util.DebugError("CodeAssessmentAdd:", err)
		return 0
	}
	return int64(data.Indices)
}

func CodeAssessmentGetByID(session *xorm.Session, id string) []model.CodeAssessmentInfo {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var ret []model.CodeAssessmentInfo
	err := session.Where("id = ?", id).Find(&ret)
	if err != nil {
		util.DebugError("CodeAssessmentAddGetByID:", err)
		return nil
	}
	return ret
}

func CodeAssessmentGetByIndices(session *xorm.Session, indices int) *model.CodeAssessmentInfo {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var ret model.CodeAssessmentInfo
	has, err := session.Where("indices = ?", indices).Get(&ret)
	if err != nil {
		util.DebugError("CodeAssessmentGetByIndices:", err)
		return nil
	}
	if !has {
		return nil
	}
	return &ret
}
