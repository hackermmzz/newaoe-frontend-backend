package Home

import (
	"encoding/json"
	"newaoe/dao"
	"newaoe/service/Code"

	"newaoe/util"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 提交历史记录（发送给前端的）
type SubmitRecord struct {
	Indices     int                `json:"indices"`
	SubmitTime  time.Time          `json:"submittime"`
	Header      string             `json:"header"`
	Source      string             `json:"source"`
	HeaderSize  int64              `json:"headersize"`
	SourceSize  int64              `json:"sourcesize"`
	Description string             `json:"description"`
	Status      Code.CodeRunStatus `json:"status"`
}

func StudentHistoryGet(ctx *gin.Context) {
	//获取用户数据
	data := util.GetCtxTookenInfo(ctx)
	if data == nil {
		util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
		return
	}
	//获取范围
	range_ := ctx.Query("range")
	parts := strings.Split(range_, ":")
	if len(parts) != 2 {
		util.ResponseNAK_MSG(ctx, "参数传递错误!", "")
		return
	}
	beg, err0 := strconv.Atoi(parts[0])
	end, err1 := strconv.Atoi(parts[1])
	if err0 != nil || err1 != nil {
		util.ResponseNAK_MSG(ctx, "参数传递错误!", "")
		return
	}
	//获取记录条数
	recordCnt := dao.CodeRunCountById(data["id"])
	//从数据库获取提交历史
	historyRecords := dao.CodeRunGetRangeById(data["id"], beg, end)
	if historyRecords == nil {
		historyRecords = make([]dao.CodeRunInfo, 0)
	}
	//编辑数据
	submitRecords := make([]SubmitRecord, len(historyRecords))
	for i := 0; i < len(submitRecords); i += 1 {
		record := &submitRecords[i]
		history := historyRecords[i]
		//
		headerInfo := dao.OssCheckFileExist(history.Header)
		sourceInfo := dao.OssCheckFileExist(history.Source)
		if headerInfo == nil || sourceInfo == nil {
			util.Debug("怎么可能出现文件不存在的情况呢?")
			continue
		}
		//
		record.Header = history.Header
		record.Source = history.Source
		record.HeaderSize = headerInfo.Size
		record.SourceSize = sourceInfo.Size
		record.Indices = history.Indices
		record.Description = history.Description
		record.SubmitTime = history.SubmitTime
		err := json.Unmarshal([]byte(history.Status), &record.Status)
		// 解析状态失败，跳过
		if err != nil {
			util.Debug("StudentHistoryGet: unmarshal status failed:" + err.Error())
			continue
		}
	}
	//
	dataSubmit := map[string]interface{}{
		"record":      submitRecords,
		"totalrecord": recordCnt,
	}
	util.ResponseACK_MSG(ctx, "历史记录获取成功", dataSubmit)
}
