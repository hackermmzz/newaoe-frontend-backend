package dao

import (
	"newaoe/Src/codeCommonSubmit/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func CodeCommonInfoAdd(session *xorm.Session, data model.CodeCommonInfo) int64 {
	_, err := session.Insert(&data)
	if err != nil {
		util.DebugError("CodeCommonInfoAdd:", err)
		return 0
	}
	return int64(data.Indices)
}

func CodeCommonGetByIndices(indices int) *model.CodeCommonInfo {
	var ret model.CodeCommonInfo
	has, err := database.DB.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		util.DebugError("CodeCommonGetByIndices:", err, indices)
		return nil
	}
	return &ret
}
func CodeCommonGetByID(id string) []model.CodeCommonInfo {
	var ret []model.CodeCommonInfo
	err := database.DB.Where("id = ?", id).Find(&ret)
	if err != nil {
		util.DebugError("CodeCommonGetByID:", err)
		return nil
	}
	return ret
}
