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
	"time"

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
	dataArr := make([]model.CodeRunInfo, len(distinctIndices))
	idx := 0
	version := util.UTC_Time().UnixNano()
	statusIndicesMap := make(map[int]codeRunStatusInfoPushRedis)
	for ind := range distinctIndices {
		data, exist := redis.RedisGet(ctx, fmt.Sprintf("CodeRun:%v", ind))
		dt := &dataArr[idx]
		dt.Indices = ind
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
			dt.Status = model.ProcessDataMessageByStatus(status).Marshal()
		}
		idx += 1
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
	updateRank(session, dataArr)
	//移除已经结束的记录
	removeCodeRunRecordAlreadyFinish(session, dataArr)
	//提交事务
	err = session.Commit()
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: commit transaction failed:", err)
		return
	}
}

func updateRank(session *xorm.Session, dt []model.CodeRunInfo) {
	//
	type CodeRunMsg struct {
		Score int  `json:"score"`
		Frame int  `json:"frame"`
		Win   bool `json:"win"`
	}
	//
	if len(dt) == 0 {
		return
	}
	//
	RankInfoArr := make([]model.RankInfo, 0)
	var err error
	var info model.CodeRunInfo
	var rankinfo model.RankInfo
	for i := range dt {
		ind := dt[i].Indices
		status := dt[i].Status //这个就是CodeRunStatusInfo
		key := fmt.Sprintf("CodeRunInfoForRank:%v", ind)
		//反序列化数据
		var codeRunStatusInfo model.CodeRunStatusInfo
		err = codeRunStatusInfo.Unmarshal([]byte(status)) //redis的数据肯定对
		if err != nil {
			util.DebugError("Why codeRunStatusInfo.Unmarshal([]byte(status)) fail?")
			continue
		}
		//检测状态码，如果不是在运行就不管
		if codeRunStatusInfo.Status < model.Code_Status_Running {
			continue
		}
		//获取当前状态
		data, exist := redis.RedisGet(context.Background(), key)
		if !exist {
			//从数据库读取
			infoPtr := dao.CodeRunGetByIndices(nil, ind)
			if infoPtr == nil {
				util.DebugError("为什么这里是空指针!")
				continue
			}
			info = *infoPtr
			//写入redis
			data, err = json.Marshal(info)
			redis.RedisSet(context.Background(), key, string(data), time.Duration(30)*time.Minute)
		}
		err = json.Unmarshal(data, &info)
		if err != nil {
			util.DebugError("updateRank: json.Unmarshal failed:" + err.Error())
			continue
		}
		//最终数据
		rankinfo.ID = info.ID
		rankinfo.SubmitTime = info.SubmitTime
		if len(info.Description) == 0 {
			info.Description = "原神启动!"
		}
		msgMp := make(map[string]string)
		msgMp["desc"] = info.Description
		msgMp["status"] = codeRunStatusInfo.Data
		msgByte, _ := json.Marshal(msgMp)
		rankinfo.Msg = string(msgByte)
		rankinfo.Score = codeRunStatusInfo.Score
		rankinfo.Frame = codeRunStatusInfo.Frame
		rankinfo.Win = codeRunStatusInfo.Win
		RankInfoArr = append(RankInfoArr, rankinfo)
	}
	//去重
	distinctID := make(map[string]model.RankInfo)
	finalRankInfo := make([]model.RankInfo, 0)
	for i := range RankInfoArr {
		info := RankInfoArr[i]
		//失败不记录到榜单
		if !info.Win {
			continue
		}
		//
		old, ok := distinctID[info.ID]
		if !ok || dao.RankIsBetter(&info, &old) {
			distinctID[info.ID] = info
		}
	}
	for _, value := range distinctID {
		finalRankInfo = append(finalRankInfo, value)
	}
	//提交更新
	if !dao.RankBatchUpdateOrInsertIfBetter(session, finalRankInfo) {
		util.DebugError("updateRank fail!")
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
	if !dao.CodeRunningBatchRemove(session, needRemove) {
		util.DebugError("批量移除CodeRunning失败!执行单个单个删除操作!")
		//单个单个移除
		for _, ind := range needRemove {
			if !dao.CodeRunningRemove(session, ind) {
				util.DebugError("单个删除CodeRunning表记录失败!")
			}
		}
	}
	return true
}
