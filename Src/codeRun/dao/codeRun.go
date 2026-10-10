package dao

import (
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"
	"strings"

	"xorm.io/xorm"
)

func CodeRunAdd(session *xorm.Session, c model.CodeRunInfo) (int, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	_, err := session.Insert(&c)
	if err != nil {
		return 0, err
	}
	return c.Indices, nil
}

func CodeRunExist(session *xorm.Session, indices int) (bool, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	exist, err := session.Where("indices= ?", indices).Exist(&model.CodeRunInfo{})
	if err != nil {
		util.Debug("CodeRunExist err:", err)
		return false, err
	}
	return exist, nil
}

func CodeRunUpdate(session *xorm.Session, newInfo model.CodeRunInfo) (bool, error) {
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
		return false, err
	}
	return affected > 0, nil
}

func CodeRunUpdateStatusAsync(session *xorm.Session, indices int, status string) (bool, error) {
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
		return false, err
	}
	return true, nil
}

func CodeRunUpdateStatusSync(session *xorm.Session, indices int, status string) (bool, error) {
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
		return false, err
	}
	return true, nil
}

func CodeRunGetByIndices(session *xorm.Session, indices int) (*model.CodeRunInfo, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var ret model.CodeRunInfo
	has, err := session.Where("indices=?", indices).Get(&ret)
	if err != nil || !has {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return &ret, nil
}

func CodeRunBatchGetByIndices(session *xorm.Session, indices []int) ([]model.CodeRunInfo, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var ret []model.CodeRunInfo
	if len(indices) == 0 {
		return ret, nil
	}
	err := session.In("indices", indices).
		Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func CodeRunGetById(session *xorm.Session, id string) ([]model.CodeRunInfo, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var ret []model.CodeRunInfo
	err := session.Where("id = ?", id).Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func CodeRunCountById(session *xorm.Session, id string) (int64, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	cnt, err := session.Where("id = ?", id).Count(&model.CodeRunInfo{})
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// 不包含end
func CodeRunGetRangeById(session *xorm.Session, id string, beg int, end int, class ...int) ([]model.CodeRunInfo, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if beg < 0 || beg >= end {
		return nil, util.NewError("beg<0 or beg>=end!")
	}
	var ret []model.CodeRunInfo
	query := session.Where("id = ?", id)
	query = getExtraClassInCondition(query, class...)
	err := query.Desc("submittime").
		Limit(end-beg, int(beg)).
		Find(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// 获取指定范围的提交记录(不包含end)
func CodeRunGetByRange(session *xorm.Session, beg int, end int, class ...int) ([]model.CodeRunInfo, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	var data []model.CodeRunInfo
	session = getExtraClassInCondition(session, class...)
	err := session.Desc("indices").Limit(end-beg, beg).Find(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// 生成IN class的额外条件
func getExtraClassInCondition(session *xorm.Session, class ...int) *xorm.Session {
	if len(class) == 0 {
		return session
	}
	if len(class) == 1 {
		return session.Where("class = ?", class[0])
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
