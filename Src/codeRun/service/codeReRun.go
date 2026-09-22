package service

import (
	"errors"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"
)

// 由于运行出问题，重新运行以前的记录（不会创建新的历史记录）
func ReRunHistoryCode(info model.CodeRunningInfo) error {
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return errors.New("开启事务失败!" + err.Error())
	}
	defer session.Rollback()
	//判断codeRunning表里面是否已经又这个记录
	has, err := dao.CodeRunningExist(session, info.Indices)
	if err != nil {
		return err
	}
	//没有则插入，好在运行异常可以恢复
	if !has {
		if !dao.CodeRunningInsert(session, model.CodeRunningInfo{
			Indices:    info.Indices,
			SubmitTime: util.UTC_Time(),
			RunType:    info.RunType,
		}) {
			return errors.New("插入CodeRunning表失败!")
		}
	} else {
		//有则更新提交时间
		if !dao.CodeRunningUpdate(session, model.CodeRunningInfo{
			Indices:    info.Indices,
			SubmitTime: util.UTC_Time(),
			RunType:    info.RunType,
		}) {
			return errors.New("更新CodeRunning表失败!")
		}
	}
	//查询以前历史记录
	historInfo := dao.CodeRunGetByIndices(session, info.Indices)
	if historInfo == nil {
		return errors.New("查询以前历史记录失败!")
	}
	//擦除以前的历史记录状态
	if !dao.CodeRunUpdateStatusAsync(session, info.Indices, model.NewCodeRunStatusInfo().Marshal()) {
		util.DebugError("RunUserCode update status fail!")
		return errors.New("擦除以前历史状态失败!")
	}
	//提交记录
	if err = session.Commit(); err != nil {
		return errors.New("提交事务失败:" + err.Error())
	}
	//加入到队列中去
	if err = addCodeFile(CodeRunTaskInfo{
		CodeRunInfo: *historInfo,
		RunType:     info.RunType,
	}); err != nil {
		return err
	}
	return nil
}
