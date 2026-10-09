package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"

	"xorm.io/xorm"
)

func UpdateRank(session *xorm.Session, dt []model.CodeRunInfo) {
	assessmentRank := make([]model.CodeRunInfo, 0)
	commonRank := make([]model.CodeRunInfo, 0)
	//分类
	for _, v := range dt {
		switch v.Class {
		case model.Code_AssessmentSubmit:
			assessmentRank = append(assessmentRank, v)
		case model.Code_CommonSubmit, model.Code_ReRunSubmit:
			commonRank = append(commonRank, v)
		default:
			util.DebugError("UpdateRank未知类型", v)
		}
	}
	//更新普通提交排行榜
	if len(commonRank) > 0 {
		err := updateCommonSubmitRank(session, commonRank)
		if err != nil {
			util.DebugError("updateCommonSubmitRank", err)
		}
	}
	//更新考核提交排行榜
	if len(assessmentRank) > 0 {
		err := updateAssessmentSubmitRank(session, assessmentRank)
		if err != nil {
			util.DebugError("updateAssessmentSubmitRank", err)
		}
	}
}

// 更新普通提交的排行榜
func updateCommonSubmitRank(session *xorm.Session, dt []model.CodeRunInfo) error {
	//
	if len(dt) == 0 {
		return nil
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
			return errors.New("Why codeRunStatusInfo.Unmarshal([]byte(status)) fail?")
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
				return errors.New("为什么这里是空指针!")
			}
			info = *infoPtr
			//写入redis
			data, err = json.Marshal(info)
			redis.RedisSet(context.Background(), key, string(data), time.Duration(30)*time.Minute)
		}
		err = json.Unmarshal(data, &info)
		if err != nil {
			return errors.New("updateRank: json.Unmarshal failed:" + err.Error())
		}
		//最终数据
		rankinfo.ID = info.ID
		rankinfo.SubmitTime = info.SubmitTime
		if len(info.Description) == 0 {
			info.Description = "原神启动!"
		}
		rankinfo.Msg = model.RankMsgInfo{
			Description: util.TruncateString(info.Description, 32*1024), //最多32kb
			Status:      codeRunStatusInfo,
		}.Marshal()
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
		//
		finalRankInfo = append(finalRankInfo, value)
	}
	//提交更新
	if !dao.RankBatchUpdateOrInsertIfBetter(session, finalRankInfo) {
		return util.NewError("updateRank fail!")
	}
	return nil
}

// 更新考核提交排行榜
func updateAssessmentSubmitRank(session *xorm.Session, dt []model.CodeRunInfo) error {
	//
	if len(dt) == 0 {
		return nil
	}
	//
	RankInfoArr := make([]model.AssessmentRankInfo, 0)
	var err error
	var info model.CodeRunInfo
	var rankinfo model.AssessmentRankInfo
	for i := range dt {
		ind := dt[i].Indices
		status := dt[i].Status //这个就是CodeRunStatusInfo
		key := fmt.Sprintf("CodeRunInfoForRank:%v", ind)
		//反序列化数据
		var codeRunStatusInfo model.CodeRunStatusInfo
		err = codeRunStatusInfo.Unmarshal([]byte(status)) //redis的数据肯定对
		if err != nil {
			return errors.New("Why codeRunStatusInfo.Unmarshal([]byte(status)) fail?")
		}
		//获取当前状态
		data, exist := redis.RedisGet(context.Background(), key)
		if !exist {
			//从数据库读取
			infoPtr := dao.CodeRunGetByIndices(nil, ind)
			if infoPtr == nil {
				return errors.New("为什么这里是空指针!")
			}
			info = *infoPtr
			//写入redis
			data, err = json.Marshal(info)
			redis.RedisSet(context.Background(), key, string(data), time.Duration(30)*time.Minute)
		}
		err = json.Unmarshal(data, &info)
		if err != nil {
			return errors.New("updateAssessmentSubmitRank: json.Unmarshal failed:" + err.Error())
		}
		//最终数据
		rankinfo.ID = info.ID
		rankinfo.SubmitTime = info.SubmitTime
		if len(info.Description) == 0 {
			info.Description = "原神启动!"
		}
		rankinfo.Msg = model.AssessmentRankMsgInfo{
			Status: codeRunStatusInfo,
		}.Marshal()
		rankinfo.Score = codeRunStatusInfo.Score
		rankinfo.Frame = codeRunStatusInfo.Frame
		rankinfo.Status = mapCodeRunStatusToAssessmenrRankStatus(codeRunStatusInfo.Status)
		RankInfoArr = append(RankInfoArr, rankinfo)
	}
	//去重
	distinctID := make(map[string]model.AssessmentRankInfo)
	finalRankInfo := make([]model.AssessmentRankInfo, 0)
	for i := range RankInfoArr {
		info := RankInfoArr[i]
		//
		old, ok := distinctID[info.ID]
		if !ok || dao.AssessmentRankIsBetter(&info, &old) {
			distinctID[info.ID] = info
		}
	}
	for _, value := range distinctID {
		//
		finalRankInfo = append(finalRankInfo, value)
	}
	//提交更新
	if !dao.AssessmentRankBatchUpdateOrInsertIfBetter(session, finalRankInfo) {
		return util.NewError("updateAssessmentSubmitRank fail!")
	}
	return nil
}

func mapCodeRunStatusToAssessmenrRankStatus(status int) int {
	if status == model.Code_Status_Compile_Fail {
		return model.AssessmenrRankStatus_CompileFail
	}
	if status == model.Code_Status_Crash {
		return model.AssessmenrRankStatus_Crash
	}
	if status == model.Code_Status_Success {
		return model.AssessmenrRankStatus_Win
	}
	if status == model.Code_Status_Fail {
		return model.AssessmenrRankStatus_Fail
	}
	return model.AssessmenrRankStatus_Error
}
