package vip

import (
	"fmt"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"
	"path"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// 定义数据结构体
type finalDataInfo struct {
	ID           string `json:"id"`            //学生学号
	Submit       bool   `json:"submit"`        //是否提交
	CompileError bool   `json:"compile_error"` //编译错误
	Crash        bool   `json:"crash"`         //崩溃
	Win          bool   `json:"win"`           //是否胜利
	Frame        int    `json:"frame"`         //运行的帧数
	Score        int    `json:"score"`         //运行的得分
}

func StudentInfoExport(ctx *gin.Context) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	//获取所有学生
	session := dao.DB.NewSession()
	defer session.Close()
	students := dao.UserGetAll(session)
	if len(students) == 0 {
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//获取所有学生最后一次提交记录
	records := dao.CodeRunGetStudentLastSubmitRecord(session)
	if len(records) == 0 {
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//建立map一一对应
	infoMap := make(map[string]dao.CodeRunInfo)
	for _, v := range records {
		infoMap[v.ID] = v
	}

	//获取所有key
	keys := util.GetJsonKeys(finalDataInfo{})
	// 写表头
	for i, h := range keys {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	// 写数据
	datas := make([]finalDataInfo, len(students))
	for i, v := range students {
		info, exist := infoMap[v.Id]
		if !exist {
			info.Status = ""
		}
		data := processFinalDataInfo(students[i].Id, info.Status)
		data.ID = students[i].Id
		datas[i] = data
	}
	//排序
	sortFinalDataArray(datas)
	//写入excel
	for i, item := range datas {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), item.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.Submit)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.CompileError)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Win)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Frame)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.Score)
	}
	//获取excl数据
	buf, err := f.WriteToBuffer()
	if err != nil {
		util.DebugError("获取Excel数据失败!", err)
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//数据存入oss
	filename := fmt.Sprintf("studentInfoExport_%v.xlsl", util.RandomString(10))
	filePath := path.Join(config.Conf.OSS.TmpBaseFolder, filename)
	if !dao.OssUploadFileData(filePath, buf.Bytes(), "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") {
		util.DebugError("保存到oss失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//获取下载链接
	expire_time := time.Duration(24) * time.Hour
	url := dao.OssGetDownloadFileUrl(filePath, expire_time, true)
	if url == "" {
		util.DebugError("获取下载链接失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//返回路径
	util.ResponseACK_MSG(ctx, "导出成功!", url)
}

func processFinalDataInfo(id string, status string) finalDataInfo {
	ret := finalDataInfo{
		Submit:       false,
		CompileError: true,
		Win:          false,
		Crash:        true,
		Score:        0,
		Frame:        0,
	}
	//如果不存在记录，
	if status == "" {
		ret.Submit = false
		return ret
	}
	ret.Submit = true
	//拿到运行状态
	var runstatus dao.CodeRunStatusInfo
	if err := runstatus.Unmarshal([]byte(status)); err != nil {
		util.DebugError("StudentInfoExport 反序列化数据错误!", err, id)
		//就默认学生是编译错误
		ret.CompileError = true
		return ret
	}
	ret.CompileError = false
	////////根据运行状态判断
	//编译失败/服务器异常/正在编译/等待
	switch runstatus.Status {
	case dao.Code_Status_Error, dao.Code_Status_Wait,
		dao.Code_Status_Compile, dao.Code_Status_Compile_Fail:
		//再单独处理一下，提醒一下
		switch runstatus.Status {
		case dao.Code_Status_Error:
			util.DebugError("服务器异常导致的错误!", id)
		case dao.Code_Status_Wait:
			util.DebugError("代码正在等待运行!", id)
		case dao.Code_Status_Compile:
			util.DebugError("代码正在编译!", id)
		}
		ret.CompileError = true
		return ret
	}
	ret.CompileError = false
	//成功/失败/崩溃 都需要得分和帧数
	ret.Score = runstatus.Score
	ret.Frame = runstatus.Frame
	//运行失败/正在运行/崩溃
	switch runstatus.Status {
	case dao.Code_Status_Running, dao.Code_Status_Fail,
		dao.Code_Status_Crash, dao.Code_Status_Compile_Success:
		switch runstatus.Status {
		case dao.Code_Status_Running:
			util.DebugError("代码正在运行!", id)
		case dao.Code_Status_Compile_Success:
			util.DebugError("代码处于编译成功状态但没推进!", id)
		}
		ret.Win = false
		ret.Crash = runstatus.Status == dao.Code_Status_Crash
		return ret
	}
	//运行成功
	ret.Win = true
	return ret
}

func sortFinalDataArray(datas []finalDataInfo) {
	sort.Slice(datas, func(i, j int) bool {
		data0 := datas[i]
		data1 := datas[j]
		return !cmpTwoFinalData(data0, data1) //逆序排序
	})
}

func cmpTwoFinalData(data0 finalDataInfo, data1 finalDataInfo) bool {
	//未提交
	if data0.Submit == false {
		return true
	}
	//编译错误
	if data0.CompileError {
		return true
	}
	//2个都崩溃
	if data0.Crash && data1.Crash {
		return data0.Score < data1.Score
	}
	//1个崩溃
	if data0.Crash || data1.Crash {
		return data0.Crash
	}
	//2个都运行失败
	if !data0.Win && !data1.Win {
		return data0.Score < data1.Score
	}
	//只有一个运行失败
	if !data0.Win || !data1.Win {
		return !data0.Win
	}
	//都运行成功
	return data0.Frame > data1.Frame
}
