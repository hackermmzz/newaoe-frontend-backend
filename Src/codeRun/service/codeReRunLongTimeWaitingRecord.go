package service

import (
	"errors"
	"newaoe/Src/codeRun/dao"
	database "newaoe/Src/databse"
	"time"
)

func ReRunLongTimeWaitRecord(expireDuration time.Duration) (int, error) {
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return 0, errors.New("开始事务失败：" + err.Error())
	}
	defer session.Rollback()
	// 获取超时的记录
	runningList := dao.CodeRunningGetExpireTime(session, expireDuration, 20)
	if len(runningList) == 0 {
		return 0, nil
	}
	//遍历每个info
	for i, info := range runningList {
		if err := ReRunHistoryCode(info); err != nil {
			return i, errors.New("重新运行失败:" + err.Error())
		}
	}
	return 0, nil
}
