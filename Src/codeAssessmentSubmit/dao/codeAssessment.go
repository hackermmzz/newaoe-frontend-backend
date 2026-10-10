package dao

import (
	"newaoe/Src/codeAssessmentSubmit/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func CodeAssessmentAdd(session *xorm.Session, data model.CodeAssessmentInfo) (int64, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	_, err := session.Insert(&data)
	if err != nil {
		return 0, err
	}
	return int64(data.Indices), nil
}

func CodeAssessmentGetByID(session *xorm.Session, id string) ([]model.CodeAssessmentInfo, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var ret []model.CodeAssessmentInfo
	err := session.Where("id = ?", id).Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func CodeAssessmentGetByIndices(session *xorm.Session, indices int) (*model.CodeAssessmentInfo, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var ret model.CodeAssessmentInfo
	has, err := session.Where("indices = ?", indices).Get(&ret)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &ret, nil
}

// 不包含end
func CodeAssessmentGetByRange(session *xorm.Session, beg int, end int) ([]model.CodeAssessmentInfo, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	if beg < 0 || beg >= end {
		return nil, util.NewError("beg<0 or beg>=end!")
	}

	var ret []model.CodeAssessmentInfo
	err := session.
		Limit(end-beg, beg).
		Find(&ret)
	if err != nil {
		return nil, util.NewError("CodeAssessmentGetByRange:", err)
	}
	return ret, nil
}
