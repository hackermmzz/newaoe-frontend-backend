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

	"xorm.io/xorm"
)

// 批量更新代码运行状态
func BatchUpdateCodeRunStatus(indices []int) {
	// 1. 对 indices 去重
	distinctIndices := make(map[int]bool)
	for _, ind := range indices {
		distinctIndices[ind] = true
	}
	if len(distinctIndices) == 0 {
		return
	}
	// 2. 从 Redis 获取数据
	ctx := context.Background()
	dataArr := make([]model.CodeRunInfo, 0)
	version := util.UTC_Time().UnixNano()
	statusIndicesMap := make(map[int]codeRunStatusInfoPushRedis)
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
			}
			dt.Status = v.CodeRunStatusInfo.Marshal()
			statusIndicesMap[ind] = v
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
		return
	}
	//创建临时表
	temp_table := fmt.Sprintf("tmp_code_run_%v", version)
	_, err = session.Exec(fmt.Sprintf("CREATE TEMPORARY TABLE %v (indices INT, status TEXT, new_version BIGINT)", temp_table))
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: create temp table failed:" + err.Error())
		session.Rollback()
		return
	}
	// 批量insert临时表
	type TempCodeRunInfo struct {
		Indices    int    `xorm:"indices"`
		Status     string `xorm:"text status"`
		NewVersion int64  `xorm:"new_version"`
	}
	tempDataArr := make([]TempCodeRunInfo, len(dataArr))
	for i := range dataArr {
		tempDataArr[i] = TempCodeRunInfo{
			Indices:    dataArr[i].Indices,
			NewVersion: dataArr[i].Version,
		}
		//解析status(它包含完整运行数据)
		tempDataArr[i].Status = statusIndicesMap[dataArr[i].Indices].CodeRunStatusInfo.Marshal()
	}
	_, err = session.Table(temp_table).Insert(tempDataArr)
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: insert temp table failed:" + err.Error())
		session.Rollback()
		return
	}
	// JOIN更新主表
	cmd := fmt.Sprintf(`	
    UPDATE %v c
    JOIN %v t ON c.indices = t.indices
    SET c.status = t.status, c.version = t.new_version
    where c.version <t.new_version
	`, model.CodeRunInfo{}.TableName(), temp_table)
	_, err = session.Exec(cmd)
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: update main table failed:" + err.Error())
		session.Rollback()
		return
	}
	//删除临时表
	_, err = session.Exec(fmt.Sprintf("DROP TEMPORARY TABLE %v", temp_table))
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: drop temp table failed:" + err.Error())
		return
	}
	//更新排行榜
	UpdateRank(session, dataArr)
	//移除已经结束的记录
	removeCodeRunRecordAlreadyFinish(session, dataArr)
	//提交事务
	err = session.Commit()
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: commit transaction failed:", err)
		return
	}
}

func removeCodeRunRecordAlreadyFinish(session *xorm.Session, dt []model.CodeRunInfo) bool {
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
	removed, err := dao.CodeRunningBatchRemove(session, needRemove)
	if err != nil {
		util.DebugError("CodeRunningBatchRemove:", err)
	}
	if !removed {
		util.DebugError("批量移除CodeRunning失败!执行单个单个删除操作!")
		//单个单个移除
		for _, ind := range needRemove {
			removed, err := dao.CodeRunningRemove(session, ind)
			if err != nil {
				util.DebugError("CodeRunningRemove:", err)
			}
			if !removed {
				util.DebugError("单个删除CodeRunning表记录失败!")
			}
		}
	}
	return true
}
