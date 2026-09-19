package service

import (
	"errors"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"time"
)

func ReRunLongTimeWaitRecord() error {
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return errors.New("开始事务失败：" + err.Error())
	}
	defer session.Rollback()
	// 获取超时的记录
	expireDuration := time.Duration(config.Conf.Code.CodeWaitTooLongTimeLimit) * time.Minute
	runningList := dao.CodeRunningGetExpireTime(session, expireDuration, 20)
	if len(runningList) == 0 {
		return nil
	}
	//遍历每个info
	for _, info := range runningList {
		if err := ReRunHistoryCode(info.Indices); err != nil {
			return errors.New("重新运行失败:" + err.Error())
		}
	}
	return nil
}
