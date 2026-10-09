package service

import (
	"errors"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/service"
)

func ReRunAnomalRecord() (int, error) {
	totalCnt := 0
	errs := make([]error, 0)
	step := 50
	for i := 0; ; i += step {
		infos, err := dao.CodeRunningGetByRange(nil, i, i+step)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if len(infos) == 0 {
			break
		}

		totalCnt += len(infos)
		//逐条运行
		for _, v := range infos {
			err = service.ReRunHistoryCode(v)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	//
	if len(errs) != 0 {
		return totalCnt, errors.Join(errs...)
	}
	return totalCnt, nil
}
