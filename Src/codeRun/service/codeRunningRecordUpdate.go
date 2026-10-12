package service

import (
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
)

func UpdateCodeRunningRecord(dt []model.CodeRunInfo) error {
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	defer session.Rollback()
	//生成数据
	codeRunningInfo := make([]model.CodeRunningInfo, 0)
	for _, info := range dt {
		runinfo := model.CodeRunStatusInfo{}
		err := runinfo.Unmarshal([]byte(info.Status))
		status := runinfo.Status
		if err != nil {
			status = model.Code_Status_Error
		}
		codeRunningInfo = append(codeRunningInfo, model.CodeRunningInfo{
			Indices: info.Indices,
			Status:  status,
		})
	}
	//更新运行记录的状态
	var better []model.CodeRunningInfo
	var err error
	if better, err = dao.CodeRunningBatchUpdateIfBetter(session, codeRunningInfo); err != nil {
		return err
	}
	//删除已完成的记录
	if err := removeRecordFinish(session, better); err != nil {
		return err
	}
	//提交事务
	if err := session.Commit(); err != nil {
		return err
	}
	return nil
}

func removeRecordFinish(session *database.Session, dt []model.CodeRunningInfo) error {
	needRemove := make([]int, 0)
	var runinfo model.CodeRunStatusInfo
	for _, info := range dt {
		//如果是已完成的记录，添加到删除列表
		if runinfo.StatusIsFinish(info.Status) {
			needRemove = append(needRemove, info.Indices)
		}
	}
	//批量删除
	_, err := dao.CodeRunningBatchRemove(session, needRemove)
	if err != nil {
		return err
	}
	return nil
}
