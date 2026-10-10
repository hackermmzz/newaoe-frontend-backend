package service

import (
	"errors"
	"fmt"
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
		inserted, err := dao.CodeRunningInsert(session, model.CodeRunningInfo{
			Indices:    info.Indices,
			SubmitTime: util.UTC_Time(),
			RunType:    info.RunType,
		})
		if err != nil {
			return err
		}
		if !inserted {
			return errors.New("插入CodeRunning表失败!")
		}
	} else {
		newInfo := model.CodeRunningInfo{
			Indices:    info.Indices,
			SubmitTime: util.UTC_Time(),
			RunType:    info.RunType,
		}
		//有则更新提交时间
		if ok, err := dao.CodeRunningUpdate(session, newInfo); err != nil || !ok {
			if err != nil {
				return err
			}
			return util.NewError(fmt.Sprintf("更新CodeRunning表失败!oldInfo:%v ,newInfo:%v , Err:%v", info, newInfo, err))
		}
	}
	//查询以前历史记录
	historInfo, err := dao.CodeRunGetByIndices(session, info.Indices)
	if err != nil {
		return err
	}
	if historInfo == nil {
		return errors.New("查询以前历史记录失败!")
	}
	//擦除以前的历史记录状态
	updated, err := dao.CodeRunUpdateStatusAsync(session, info.Indices, model.NewCodeRunStatusInfo().Marshal())
	if err != nil {
		return err
	}
	if !updated {
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
	}, 3); err != nil {
		return err
	}
	return nil
}
