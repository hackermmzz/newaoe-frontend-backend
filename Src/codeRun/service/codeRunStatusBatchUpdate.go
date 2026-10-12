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
	if err = dao.CodeRunBatchUpdateStatusIfBetter(session, dataArr); err != nil {
		util.DebugError("CodeRunBatchUpdateStatus:", err)
		session.Rollback()
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
