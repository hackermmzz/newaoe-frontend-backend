package dao

import (
	"fmt"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func CodeRunAdd(session *xorm.Session, c model.CodeRunInfo) int {
	_, err := session.Insert(&c)
	if err != nil {
		util.DebugError("CodeRunAdd:", err)
		return 0
	}
	return c.Indices
}

func CodeRunExist(session *xorm.Session, indices int) bool {
	exist, err := session.Where("indices= ?", indices).Exist(&model.CodeRunInfo{})
	if err != nil {
		util.Debug("CodeRunExist err:", err)
		return false
	}
	return exist
}

func CodeRunUpdate(session *xorm.Session, newInfo model.CodeRunInfo) bool {
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

func CodeRunUpdateStatusSync(indices int, status string) bool {
	var info model.CodeRunInfo
	info.Status = status
	_, err := database.DB.Where("indices=?", indices).
		Update(info)
	if err != nil {
		util.DebugError("CodeRunUpdateStatusSync:", err)
		return false
	}
	return true
}

func CodeRunGetByIndices(indices int) *model.CodeRunInfo {
	var ret model.CodeRunInfo
	has, err := database.DB.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		util.DebugError("CodeRunGetByIndices:", err)
		return nil
	}
	return &ret
}

func CodeRunBatchGetByIndices(indices []int) []model.CodeRunInfo {
	var ret []model.CodeRunInfo
	if len(indices) == 0 {
		return ret
	}
	err := database.DB.In("indices", indices).
		Find(&ret)
	if err != nil {
		util.DebugError("CodeRunBatchGetByIndices:", err)
		return nil
	}
	return ret
}

func CodeRunGetById(id string) []model.CodeRunInfo {
	var ret []model.CodeRunInfo
	err := database.DB.Where("id = ?", id).Find(&ret)
	if err != nil {
		util.DebugError("CodeRunGetById:", err)
		return nil
	}
	return ret
}

func CodeRunCountById(id string) int64 {
	cnt, err := database.DB.Where("id = ?", id).Count(&model.CodeRunInfo{})
	if err != nil {
		util.DebugError("CodeRunCountById:", err)
		return 0
	}
	return cnt
}

// 不包含end
func CodeRunGetRangeById(id string, beg int, end int) []model.CodeRunInfo {
	if beg < 0 || beg >= end {
		return nil
	}
	var ret []model.CodeRunInfo
	err := database.DB.Where("id = ?", id).
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
