package service

import (
	"errors"
	"newaoe/Src/codeRun/dao"
	"time"
)

func ReRunLongTimeWaitRecord(expireDuration time.Duration) (int, error) {
	// 获取超时的记录
	runningList := dao.CodeRunningGetExpireTime(nil, expireDuration, 20)
	if len(runningList) == 0 {
		return 0, nil
	}
	//遍历每个info
	errs := make([]error, 0)
	sucessCnt := 0
	for _, info := range runningList {
		if err := ReRunHistoryCode(info); err != nil {
			errs = append(errs, errors.New("重新运行失败:"+err.Error()))
		} else {
			sucessCnt += 1
		}
	}
	return sucessCnt, errors.Join(errs...)
}
