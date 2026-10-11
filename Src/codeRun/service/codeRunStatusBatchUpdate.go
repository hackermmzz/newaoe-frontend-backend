package service

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
	"newaoe/Src/util"
)

// 批量更新代码运行状态
func BatchUpdateCodeRunStatus(indices []int) ([]model.CodeRunInfo, error) {
	// 1. 对 indices 去重
	distinctIndices := make(map[int]bool)
	for _, ind := range indices {
		distinctIndices[ind] = true
	}
	if len(distinctIndices) == 0 {
		return nil, nil
	}
	// 2. 从 Redis 获取数据
	ctx := context.Background()
	dataArr := make([]model.CodeRunInfo, 0)
	version := util.UTC_Time().UnixNano()
	for ind := range distinctIndices {
		data, exist, err := redis.RedisGet(ctx, fmt.Sprintf("CodeRun:%v", ind))
		if err != nil {
			util.DebugError("RedisGet CodeRun status:", err)
			continue
		}
		//获取数据库里面的数据
		dbRet, err := dao.CodeRunGetByIndices(nil, ind)
		if err != nil {
			util.DebugError("CodeRunGetByIndices:", err)
			continue
		}
		if dbRet == nil {
			util.DebugError("为什么获取不到这个indice的数据呢?", ind)
			continue
		}
		//处理一下
		dt := *dbRet
		dt.Version = version
		if exist {
			var v codeRunStatusInfoPushRedis
			err := json.Unmarshal(data, &v)
			if err != nil {
				util.DebugError("Why err := json.Unmarshal(data, &v) fail?", err)
				continue //丢弃这个记录
			}
			dt.Status = v.CodeRunStatusInfo.Marshal()
		} else {
			status := model.NewCodeRunStatusInfo()
			status.Status = model.Code_Status_Error
			dt.Status = util.TruncateString(model.ProcessDataMessageByStatus(status).Marshal(), 32*1024) //截断32kb
		}
		//
		dataArr = append(dataArr, dt)
	}
	// 3. 批量更新到数据库
	// 批量insert临时表
	var err error
	session := database.NewSession()
	defer session.Close()
	//开启事务
	err = session.Begin()
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus:" + err.Error())
		return nil, err
	}
	if err = dao.CodeRunBatchUpdateStatus(session, dataArr); err != nil {
		util.DebugError("CodeRunBatchUpdateStatus:", err)
		session.Rollback()
		return nil, err
	}
	//移除已经结束的记录
	err = removeCodeRunRecordAlreadyFinish(session, dataArr)
	if err != nil {
		return nil, err
	}
	//提交事务
	err = session.Commit()
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: commit transaction failed:", err)
		return nil, err
	}
	return dataArr, nil
}

func removeCodeRunRecordAlreadyFinish(session *database.Session, dt []model.CodeRunInfo) error {
	needRemove := make([]int, 0)
	for _, info := range dt {
		ind := info.Indices
		status := info.Status
		//解析运行结果
		var codeRunningStatus model.CodeRunStatusInfo
		if err := codeRunningStatus.Unmarshal([]byte(status)); err != nil {
			util.DebugError("status解析失败:", err)
			continue
		}
		//运行状态表明程序没有结束
		if !codeRunningStatus.IsFinish() {
			continue
		}
		//
		needRemove = append(needRemove, ind)
	}
	//批量删除
	_, err := dao.CodeRunningBatchRemove(session, needRemove)
	if err != nil {
		return err
	}
	return nil
}
