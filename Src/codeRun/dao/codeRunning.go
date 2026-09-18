package dao

import (
	"newaoe/Src/codeRun/model"
	"newaoe/Src/util"
	"time"

	"xorm.io/xorm"
)

func CodeRunningInsert(session *xorm.Session, data model.CodeRunningInfo) bool {
	_, err := session.Insert(data)
	if err != nil {
		util.DebugError("CodeRunningInsert:", err)
		return false
	}
	return true
}

func CodeRunningUpdate(session *xorm.Session, info model.CodeRunningInfo) bool {
	affected, err := session.ID(info.Indices).Update(info)
	if err != nil {
		util.DebugError("CodeRunningUpdate:", err)
		return false
	}
	return affected > 0
}

func CodeRunningRemove(session *xorm.Session, indices int) bool {
	d := &model.CodeRunningInfo{
		Indices: indices,
	}
	_, err := session.Delete(d)
	if err != nil {
		util.DebugError("CodeRunningRemove:", err)
		return false
	}
	return true
}

func CodeRunningBatchRemove(session *xorm.Session, indices []int) bool {
	if len(indices) == 0 {
		return true
	}
	//
	_, err := session.In("indices", indices).Delete(&model.CodeRunningInfo{})
	if err != nil {
		util.DebugError("CodeRunningBatchRemove:", err)
		return false
	}
	return true
}

func CodeRunningGetExpireTime(session *xorm.Session, expireDuration time.Duration, number int) []model.CodeRunningInfo {
	var result []model.CodeRunningInfo
	// 当前时间减去超时时间
	expireTime := time.Now().Add(-expireDuration)
	err := session.
		Where("submittime <= ?", expireTime).
		Limit(number). // 最多返回 number 条
		Find(&result)
	if err != nil {
		util.Debug("CodeRunningGetExpireTime:", err)
		return nil
	}
	return result
}

func CodeRunningExist(session *xorm.Session, indices int) bool {
	exists, err := session.Where("indices = ?", indices).Exist(&model.CodeRunningInfo{})
	if err != nil {
		util.Debug("CodeRunningExist err:", err)
		return false
	}
	return exists
}
