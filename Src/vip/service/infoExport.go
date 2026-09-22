package service

import (
	"errors"
	"fmt"
	CodeRunDao "newaoe/Src/codeRun/dao"
	CodeRunModel "newaoe/Src/codeRun/model"
	"newaoe/Src/common/upload"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/oss"
	UserDao "newaoe/Src/user/dao"
	"newaoe/Src/util"
	"path"
	"sort"
	"time"

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

func ExcelInfoExport() (string, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	//获取所有学生
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return "", errors.New("数据库异常!")
	}
	defer session.Rollback()
	//
	students := UserDao.UserGetAll(session)
	if len(students) == 0 {
		return "", errors.New("数据库异常!")
	}
	//获取所有学生最后一次提交记录
	records := CodeRunDao.CodeRunGetStudentLastSubmitRecord(session)
	if len(records) == 0 {
		return "", errors.New("数据库异常!")
	}
	//建立map一一对应
	infoMap := make(map[string]CodeRunModel.CodeRunInfo)
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
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Crash)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Win)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.Frame)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), item.Score)
	}
	//获取excl数据
	buf, err := f.WriteToBuffer()
	if err != nil {
		return "", util.NewError("数据库异常!", err)
	}
	//数据存入oss
	filename := fmt.Sprintf("studentInfoExport_%s.xlsl", util.UUID())
	filePath := path.Join(config.Conf.OSS.TmpBaseFolder, filename)
	if !upload.UploadData(filePath, buf.Bytes(), "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet") {
		return "", util.NewError("保存到oss失败!")
	}
	//获取下载链接
	expire_time := time.Duration(24) * time.Hour
	url := oss.OssGetDownloadFileUrl(filePath, expire_time, true)
	if url == "" {
		return "", util.NewError("获取下载链接失败!", err)
	}
	//
	return url, nil
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
	var runstatus CodeRunModel.CodeRunStatusInfo
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
	case CodeRunModel.Code_Status_Error, CodeRunModel.Code_Status_Wait,
		CodeRunModel.Code_Status_Compile, CodeRunModel.Code_Status_Compile_Fail:
		//再单独处理一下，提醒一下
		switch runstatus.Status {
		case CodeRunModel.Code_Status_Error:
			util.DebugError("服务器异常导致的错误!", id)
		case CodeRunModel.Code_Status_Wait:
			util.DebugError("代码正在等待运行!", id)
		case CodeRunModel.Code_Status_Compile:
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
	case CodeRunModel.Code_Status_Running, CodeRunModel.Code_Status_Fail,
		CodeRunModel.Code_Status_Crash, CodeRunModel.Code_Status_Compile_Success:
		switch runstatus.Status {
		case CodeRunModel.Code_Status_Running:
			util.DebugError("代码正在运行!", id)
		case CodeRunModel.Code_Status_Compile_Success:
			util.DebugError("代码处于编译成功状态但没推进!", id)
		}
		ret.Win = false
		ret.Crash = runstatus.Status == CodeRunModel.Code_Status_Crash
		return ret
	}
	//运行成功
	ret.Win = true
	return ret
}
