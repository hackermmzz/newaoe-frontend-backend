package dao

import (
	"newaoe/Src/codeCommonSubmit/model"
	database "newaoe/Src/databse"

	"xorm.io/xorm"
)

func CodeCommonInfoAdd(session *xorm.Session, data model.CodeCommonInfo) (int64, error) {
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

func CodeCommonGetByIndices(session *xorm.Session, indices int) (*model.CodeCommonInfo, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var ret model.CodeCommonInfo
	has, err := session.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return &ret, nil
}
func CodeCommonGetByID(session *xorm.Session, id string) ([]model.CodeCommonInfo, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var ret []model.CodeCommonInfo
	err := session.Where("id = ?", id).Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
