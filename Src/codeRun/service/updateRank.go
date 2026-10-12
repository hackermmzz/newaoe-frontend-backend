package service

import (
	"errors"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"
)

func UpdateRank(dt []model.CodeRunInfo) error {
	//
	if len(dt) == 0 {
		return nil
	}
	session := database.NewSession()
	defer session.Close()
	//开启事务
	err := session.Begin()
	if err != nil {
		return util.NewError("UpdateRank: begin transaction failed:", err)
	}
	defer session.Rollback()
	//
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
			return util.NewError("UpdateRank未知类型", v)
		}
	}
	//更新普通提交排行榜
	if len(commonRank) > 0 {
		err := updateCommonSubmitRank(session, commonRank)
		if err != nil {
			return util.NewError("updateCommonSubmitRank", err)
		}
	}
	//更新考核提交排行榜
	if len(assessmentRank) > 0 {
		err := updateAssessmentSubmitRank(session, assessmentRank)
		if err != nil {
			return util.NewError("updateAssessmentSubmitRank", err)
		}
	}
	//提交事务
	err = session.Commit()
	if err != nil {
		return util.NewError("UpdateRank: commit transaction failed:", err)
	}
	return nil
}

// 更新普通提交的排行榜
func updateCommonSubmitRank(session *database.Session, dt []model.CodeRunInfo) error {
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
		//从数据库读取
		infoPtr, err := dao.CodeRunGetByIndices(nil, ind)
		if err != nil {
			continue //读取失败，跳过
		}
		if infoPtr == nil {
			return errors.New("为什么这里是空指针!")
		}
		info = *infoPtr
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
	updated, err := dao.RankBatchUpdateOrInsertIfBetter(session, finalRankInfo)
	if err != nil {
		return err
	}
	if !updated {
		return util.NewError("updateRank fail!")
	}
	return nil
}

// 更新考核提交排行榜
func updateAssessmentSubmitRank(session *database.Session, dt []model.CodeRunInfo) error {
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
		//反序列化数据
		var codeRunStatusInfo model.CodeRunStatusInfo
		err = codeRunStatusInfo.Unmarshal([]byte(status)) //redis的数据肯定对
		if err != nil {
			return errors.New("Why codeRunStatusInfo.Unmarshal([]byte(status)) fail?")
		}
		//获取当前状态
		//从数据库读取
		infoPtr, err := dao.CodeRunGetByIndices(nil, ind)
		if err != nil {
			continue //读取失败，跳过
		}
		if infoPtr == nil {
			return errors.New("为什么这里是空指针!")
		}
		info = *infoPtr
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
		rankinfo.Status = codeRunStatusInfo.Status
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
	err = dao.AssessmentRankBatchUpdateOrInsertIfBetter(session, finalRankInfo)
	if err != nil {
		return err
	}
	return nil
}
