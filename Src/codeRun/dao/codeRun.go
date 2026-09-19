package dao

import (
	"fmt"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"

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
func CodeRunGetRangeById(session *xorm.Session, id string, beg int, end int) []model.CodeRunInfo {
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
	err := session.Where("id = ?", id).
		OrderBy("submittime desc").
		Limit(end-beg, int(beg)).
		Find(&ret)
	if err != nil {
		util.DebugError("CodeRunGetRangeById:", err)
		return nil
	}
	return ret
}

// 获取所有学生最后一次提交记录
func CodeRunGetStudentLastSubmitRecord(session *xorm.Session) []model.CodeRunInfo {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var result []model.CodeRunInfo

	sql := fmt.Sprintf(`
	SELECT *FROM (
		SELECT *,
		ROW_NUMBER() OVER(
			PARTITION BY id 
			ORDER BY indices DESC
		) AS rn
		FROM %s
	) t
	WHERE rn = 1
	`, model.CodeRunInfo{}.TableName())
	err := session.SQL(sql).Find(&result)
	if err != nil {
		util.DebugError("CodeRunGetStudentLastSubmitRecord:", err)
		return nil
	}

	return result
}
