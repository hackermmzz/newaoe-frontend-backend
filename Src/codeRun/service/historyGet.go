package service

import (
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	"newaoe/Src/oss"
	"newaoe/Src/util"
	"time"
)

// 提交历史记录（发送给前端的）
type SubmitRecord struct {
	ID          string                  `json:"id"`
	Indices     int                     `json:"indices"`
	SubmitTime  time.Time               `json:"submittime"`
	Header      string                  `json:"header"`
	Source      string                  `json:"source"`
	HeaderSize  int64                   `json:"headersize"`
	SourceSize  int64                   `json:"sourcesize"`
	Description string                  `json:"description"`
	Status      model.CodeRunStatusInfo `json:"status"`
	Class       int                     `json:"class"`
}

// 如果id==""，表示仅按范围查找(包含end)
func GetHistoryRangeById(id string, beg int, end int, class ...int) ([]SubmitRecord, error) {
	//从数据库获取提交历史
	var historyRecords []model.CodeRunInfo
	var err error
	if id != "" {
		historyRecords, err = dao.CodeRunGetRangeById(nil, id, beg, end+1, class...)
	} else {
		historyRecords, err = dao.CodeRunGetByRange(nil, beg, end+1, class...)
	}
	if err != nil {
		return nil, err
	}
	if historyRecords == nil {
		historyRecords = make([]model.CodeRunInfo, 0)
	}
	//编辑数据
	submitRecords := make([]SubmitRecord, 0)
	//判断文件是否存在
	allfiles := make([]string, len(historyRecords)*2)
	for i, d := range historyRecords {
		allfiles[i*2] = d.Header
		allfiles[i*2+1] = d.Source
	}
	fileInfos := oss.OssCheckFilesExist(allfiles)
	//
	for _, v := range historyRecords {
		record := SubmitRecord{}
		history := v
		//
		headerInfo := fileInfos[history.Header]
		sourceInfo := fileInfos[history.Source]
		if headerInfo == nil || sourceInfo == nil {
			util.DebugError("StudentHistoryGet: file not exist: " + history.Header + " or " + history.Source)
			continue //跳过这个记录
		}
		//
		record.ID = history.ID
		record.Header = history.Header
		record.Source = history.Source
		record.HeaderSize = headerInfo.Size
		record.SourceSize = sourceInfo.Size
		record.Indices = history.Indices
		record.Description = history.Description
		record.SubmitTime = history.SubmitTime
		record.Class = history.Class
		err := record.Status.Unmarshal([]byte(history.Status))
		// 解析状态失败，跳过（兼容以前的格式)
		if err != nil {
			util.DebugError("StudentHistoryGet: unmarshal status failed:" + err.Error())
			continue
		}
		submitRecords = append(submitRecords, record)
	}
	return submitRecords, nil
}
