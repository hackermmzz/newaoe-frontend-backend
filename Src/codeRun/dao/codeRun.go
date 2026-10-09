package dao

import (
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"
	"strings"

	"xorm.io/xorm"
)

func CodeRunAdd(session *xorm.Session, c model.CodeRunInfo) int {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	_, err := session.Insert(&c)
	if err != nil {
		util.DebugError("CodeRunAdd:", err)
		return 0
	}
	return c.Indices
}

func CodeRunExist(session *xorm.Session, indices int) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	exist, err := session.Where("indices= ?", indices).Exist(&model.CodeRunInfo{})
	if err != nil {
		util.Debug("CodeRunExist err:", err)
		return false
	}
	return exist
}

func CodeRunUpdate(session *xorm.Session, newInfo model.CodeRunInfo) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	affected, err := session.
		ID(newInfo.Indices).
		AllCols().
		Update(&newInfo)
	if err != nil {
		util.DebugError("CodeRunUpdate:", err)
		return false
	}
	return affected > 0
}

func CodeRunUpdateStatusAsync(session *xorm.Session, indices int, status string) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var info model.CodeRunInfo
	info.Status = status
	_, err := session.Where("indices=?", indices).
		Update(info)
	if err != nil {
		util.DebugError("CodeRunUpdateStatusAsync:", err)
		return false
	}
	return true
}

func CodeRunUpdateStatusSync(session *xorm.Session, indices int, status string) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var info model.CodeRunInfo
	info.Status = status
	_, err := session.Where("indices=?", indices).
		Update(info)
	if err != nil {
		util.DebugError("CodeRunUpdateStatusSync:", err)
		return false
	}
	return true
}

func CodeRunGetByIndices(session *xorm.Session, indices int) *model.CodeRunInfo {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var ret model.CodeRunInfo
	has, err := session.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		util.DebugError("CodeRunGetByIndices:", err)
		return nil
	}
	return &ret
}

func CodeRunBatchGetByIndices(session *xorm.Session, indices []int) []model.CodeRunInfo {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var ret []model.CodeRunInfo
	if len(indices) == 0 {
		return ret
	}
	err := session.In("indices", indices).
		Find(&ret)
	if err != nil {
		util.DebugError("CodeRunBatchGetByIndices:", err)
		return nil
	}
	return ret
}

func CodeRunGetById(session *xorm.Session, id string) []model.CodeRunInfo {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var ret []model.CodeRunInfo
	err := session.Where("id = ?", id).Find(&ret)
	if err != nil {
		util.DebugError("CodeRunGetById:", err)
		return nil
	}
	return ret
}

func CodeRunCountById(session *xorm.Session, id string) int64 {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	cnt, err := session.Where("id = ?", id).Count(&model.CodeRunInfo{})
	if err != nil {
		util.DebugError("CodeRunCountById:", err)
		return 0
	}
	return cnt
}

// 不包含end
func CodeRunGetRangeById(session *xorm.Session, id string, beg int, end int, class ...int) []model.CodeRunInfo {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if beg < 0 || beg >= end {
		return nil
	}
	var ret []model.CodeRunInfo
	query := session.Where("id = ?", id)
	query = getExtraClassInCondition(query, class...)
	err := query.Desc("submittime").
		Limit(end-beg, int(beg)).
		Find(&ret)
	if err != nil {
		util.DebugError("CodeRunGetRangeById:", err)
		return nil
	}
	return ret
}

// 获取指定范围的提交记录(不包含end)
func CodeRunGetByRange(session *xorm.Session, beg int, end int, class ...int) []model.CodeRunInfo {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	var data []model.CodeRunInfo
	session = getExtraClassInCondition(session, class...)
	err := session.Desc("indices").Limit(end-beg, beg).Find(&data)

	if err != nil {
		util.DebugError("CodeRunGetByRange:", err)
		return nil
	}

	return data
}

// 生成IN class的额外条件
func getExtraClassInCondition(session *xorm.Session, class ...int) *xorm.Session {
	if len(class) == 0 {
		return session
	}
	placeholders := make([]string, len(class))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	cond := "class IN (" + strings.Join(placeholders, ",") + ")"
	args := make([]interface{}, len(class))
	for i, v := range class {
		args[i] = v
	}
	return session.And(cond, args...)
}
