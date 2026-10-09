package service

import (
	"errors"
	codeAssementDao "newaoe/Src/codeAssessmentSubmit/dao"
	codeRunModel "newaoe/Src/codeRun/model"
	"newaoe/Src/codeRun/service"
)

func RunAllAssessmentSubmit() (int, error) {
	totalCnt := 0
	errs := make([]error, 0)
	step := 50
	for i := 0; ; i += step {
		arr, err := codeAssementDao.CodeAssessmentGetByRange(nil, i, i+step)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if len(arr) == 0 {
			break
		}
		//运行
		for _, v := range arr {
			err := service.CodeRun(v.Indices, v.ID, codeRunModel.Code_AssessmentSubmit, codeRunModel.CodeRunType_ReleaseRun)
			if err != nil {
				errs = append(errs, err)
			} else {
				totalCnt += 1
			}
		}
	}
	//
	if len(errs) != 0 {
		return totalCnt, errors.Join(errs...)
	}
	return totalCnt, nil
}
