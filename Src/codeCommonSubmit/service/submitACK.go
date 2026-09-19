package service

import (
	"errors"
	"newaoe/Src/codeCommonSubmit/dao"
	"newaoe/Src/codeCommonSubmit/model"
	"newaoe/Src/common/upload"
	database "newaoe/Src/databse"
	"newaoe/Src/oss"
	"newaoe/Src/util"
)

func CodeCommonSubmitACK(id string, key string) (int64, error) {
	//获取文件存储路径
	filePath := upload.UploadParseAndDelKey(key)
	if len(filePath) != 3 {
		return int64(0), errors.New("这可能是个Bug")
	}
	//获取信息
	fileInfos := oss.OssCheckFilesExist(filePath)
	if len(fileInfos) != len(filePath) {
		return int64(0), errors.New("为什么长度不相等!")
	}
	headerInfo, ok0 := fileInfos[filePath[0]]
	sourceInfo, ok1 := fileInfos[filePath[1]]
	descriptionInfo, ok2 := fileInfos[filePath[2]]
	if !ok1 || !ok2 || !ok0 || descriptionInfo == nil || sourceInfo == nil || headerInfo == nil {
		return int64(0), errors.New("无法获取到文件信息")
	}
	description := oss.OssGetFileData(descriptionInfo.Key)
	//
	info := model.CodeCommonInfo{
		Indices:     0,
		ID:          id,
		Header:      headerInfo.Key,
		Source:      sourceInfo.Key,
		Description: string(description),
		HeaderSize:  headerInfo.Size,
		SourceSize:  sourceInfo.Size,
		UploadTime:  util.UTC_Time(),
	}
	//
	session := database.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		return int64(0), errors.New("开启事务失败")
	}
	defer session.Rollback()
	//写入数据库
	indices := int64(0)
	if indices = dao.CodeCommonInfoAdd(session, info); indices == 0 {
		return int64(0), errors.New("记录写入数据库失败!")
	}
	//提交事务
	if err = session.Commit(); err != nil {
		return int64(0), errors.New("提交事务失败")
	}
	//
	return indices, nil
}
