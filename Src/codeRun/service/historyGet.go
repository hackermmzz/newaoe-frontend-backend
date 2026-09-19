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

func GetHistoryRangeById(id string, beg int, end int) ([]SubmitRecord, error) {
	//从数据库获取提交历史
	historyRecords := dao.CodeRunGetRangeById(nil, id, beg, end+1)
	if historyRecords == nil {
		historyRecords = make([]model.CodeRunInfo, 0)
	}
	//编辑数据
	submitRecords := make([]SubmitRecord, len(historyRecords))
	//判断文件是否存在
	allfiles := make([]string, len(historyRecords)*2)
	for i, d := range historyRecords {
		allfiles[i*2] = d.Header
		allfiles[i*2+1] = d.Source
	}
	fileInfos := oss.OssCheckFilesExist(allfiles)
	//
	for i := 0; i < len(submitRecords); i += 1 {
		record := &submitRecords[i]
		history := historyRecords[i]
		//
		headerInfo := fileInfos[history.Header]
		sourceInfo := fileInfos[history.Source]
		if headerInfo == nil || sourceInfo == nil {
			return nil, util.NewError("怎么可能出现文件不存在的情况呢?")
		}
		//
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
	}
	return submitRecords, nil
}
