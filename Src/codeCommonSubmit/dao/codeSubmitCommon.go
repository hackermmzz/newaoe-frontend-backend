package dao

import (
	"newaoe/Src/codeCommonSubmit/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func CodeCommonInfoAdd(session *xorm.Session, data model.CodeCommonInfo) int64 {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	_, err := session.Insert(&data)
	if err != nil {
		util.DebugError("CodeCommonInfoAdd:", err)
		return 0
	}
	return int64(data.Indices)
}

func CodeCommonGetByIndices(session *xorm.Session, indices int) *model.CodeCommonInfo {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var ret model.CodeCommonInfo
	has, err := session.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		util.DebugError("CodeCommonGetByIndices:", err, indices)
		return nil
	}
	return &ret
}
func CodeCommonGetByID(session *xorm.Session, id string) []model.CodeCommonInfo {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	var ret []model.CodeCommonInfo
	err := session.Where("id = ?", id).Find(&ret)
	if err != nil {
		util.DebugError("CodeCommonGetByID:", err)
		return nil
	}
	return ret
}
