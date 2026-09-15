package Home

import (
	"newaoe/dao"

	"newaoe/util"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 提交历史记录（发送给前端的）
type SubmitRecord struct {
	Indices     int                   `json:"indices"`
	SubmitTime  time.Time             `json:"submittime"`
	Header      string                `json:"header"`
	Source      string                `json:"source"`
	HeaderSize  int64                 `json:"headersize"`
	SourceSize  int64                 `json:"sourcesize"`
	Description string                `json:"description"`
	Status      dao.CodeRunStatusInfo `json:"status"`
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
	//获取查寻的id
	student_id := ctx.Query("student_id")
	id := data["id"].(string)
	if student_id != "" {
		id = student_id
	}
	//获取历史记录
	submitRecords := getHistoryRangeById(id, beg, end)
	//
	util.ResponseACK_MSG(ctx, "历史记录获取成功", submitRecords)
}

func getHistoryRangeById(id string, beg int, end int) []SubmitRecord {
	//从数据库获取提交历史
	historyRecords := dao.CodeRunGetRangeById(id, beg, end+1)
	if historyRecords == nil {
		historyRecords = make([]dao.CodeRunInfo, 0)
	}
	//编辑数据
	submitRecords := make([]SubmitRecord, len(historyRecords))
	//判断文件是否存在
	allfiles := make([]string, len(historyRecords)*2)
	for i, d := range historyRecords {
		allfiles[i*2] = d.Header
		allfiles[i*2+1] = d.Source
	}
	fileInfos := dao.OssCheckFilesExist(allfiles)
	//
	for i := 0; i < len(submitRecords); i += 1 {
		record := &submitRecords[i]
		history := historyRecords[i]
		//
		headerInfo := fileInfos[history.Header]
		sourceInfo := fileInfos[history.Source]
		if headerInfo == nil || sourceInfo == nil {
			util.DebugError("怎么可能出现文件不存在的情况呢?")
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
		err := record.Status.Unmarshal([]byte(history.Status))
		// 解析状态失败，跳过
		if err != nil {
			util.DebugError("StudentHistoryGet: unmarshal status failed:" + err.Error())
			continue
		}
	}
	return submitRecords
}
