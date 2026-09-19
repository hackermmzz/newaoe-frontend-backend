package service

import (
	"newaoe/Src/home/dao"
	"newaoe/Src/home/model"
	"newaoe/Src/oss"
	"newaoe/Src/util"
	"path/filepath"
	"time"
)

type FeedbackRecordForManagerInfo struct {
	model.FeedbackInfo
	Link string `json:"string"`
}

func FetchFeedbackRecord(beg int, end int) ([]FeedbackRecordForManagerInfo, error) {
	//查询反馈记录
	record := dao.FeedbackGetRange(nil, beg, end)
	if record == nil {
		return nil, util.NewError("获取记录失败!")
	}
	//获取相应的html链接
	ret := make([]FeedbackRecordForManagerInfo, 0)
	for _, v := range record {
		htmlFilePath := filepath.Join(v.BaseFolder, "feedback.html")
		//获取下载链接
		expire_time := time.Duration(24) * time.Hour
		url := oss.OssGetDownloadFileUrl(htmlFilePath, expire_time, false)
		if url == "" {
			util.DebugError("获取下载链接失败!", htmlFilePath)
			continue
		}
		ret = append(ret, FeedbackRecordForManagerInfo{
			FeedbackInfo: v,
			Link:         url,
		})
	}
	//
	return ret, nil
}
