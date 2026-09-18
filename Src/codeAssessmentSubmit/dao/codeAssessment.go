package dao

import (
	"newaoe/Src/codeAssessmentSubmit/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func CodeAssessmentAdd(session *xorm.Session, data model.CodeAssessmentInfo) int64 {
	_, err := session.Insert(&data)
	if err != nil {
		util.DebugError("CodeAssessmentAdd:", err)
		return 0
	}
	return int64(data.Indices)
}

func CodeAssessmentGetByID(id string) []model.CodeAssessmentInfo {
	var ret []model.CodeAssessmentInfo
	err := database.DB.Where("id = ?", id).Find(&ret)
	if err != nil {
		util.DebugError("CodeAssessmentAddGetByID:", err)
		return nil
	}
	return ret
}

func CodeAssessmentGetByIndices(indices int) *model.CodeAssessmentInfo {
	var ret model.CodeAssessmentInfo
	has, err := database.DB.Where("indices = ?", indices).Get(&ret)
	if err != nil {
		util.DebugError("CodeAssessmentGetByIndices:", err)
		return nil
	}
	if !has {
		return nil
	}
	return &ret
}
