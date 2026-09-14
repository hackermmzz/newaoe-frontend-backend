package dao

import (
	"newaoe/util"
	"time"

	"xorm.io/xorm"
)

// 这个表负责维护哪些正在跑，但是可能会出问题的运行记录（保证失败了可以重新推入队列）
type CodeRunningInfo struct {
	Indices    int       `json:"indices" xorm:"indices pk "`
	SubmitTime time.Time `json:"submittime" xorm:"submittime"` //提交日期
}

func (CodeRunningInfo) TableName() string {
	return "CodeRunning" // 这里返回你想要的表名
}

func CodeRunningInsert(session *xorm.Session, data CodeRunningInfo) bool {
	_, err := session.Insert(data)
	if err != nil {
		util.DebugError("CodeRunningInsert:", err)
		return false
	}
	return true
}

func CodeRunningUpdate(session *xorm.Session, info CodeRunningInfo) bool {
	affected, err := session.ID(info.Indices).Update(info)
	if err != nil {
		util.DebugError("CodeRunningUpdate:", err)
		return false
	}
	return affected > 0
}

func CodeRunningRemove(session *xorm.Session, indices int) bool {
	d := &CodeRunningInfo{
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
	d := make([]CodeRunningInfo, len(indices))
	for i, v := range indices {
		d[i].Indices = v
	}
	//
	_, err := session.In("indices", d).Delete(&CodeRunningInfo{})
	if err != nil {
		util.DebugError("CodeRunningBatchRemove:", err)
		return false
	}
	return true
}

func CodeRunningGetExpireTime(session *xorm.Session, expireDuration time.Duration, number int) []CodeRunningInfo {
	var result []CodeRunningInfo
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
	exists, err := session.Where("indices = ?", indices).Exist(&CodeRunningInfo{})
	if err != nil {
		util.Debug("CodeRunningExist err:", err)
		return false
	}
	return exists
}
