package service

import (
	"errors"
	"newaoe/Src/codeAssessmentSubmit/dao"
	"newaoe/Src/codeAssessmentSubmit/model"
	"newaoe/Src/common/upload"
	database "newaoe/Src/databse"
	"newaoe/Src/oss"
	"newaoe/Src/util"
)

// 返回插入的数据库唯一键
func AssessmentSubmitACK(id string, key string, teacherSelected string) (int64, error) {
	//判断教师存在不存在
	if !dao.TeacherExist(nil, teacherSelected) {
		return int64(0), errors.New("教师" + teacherSelected + "不存在!")
	}
	//获取对应的文件路径
	filePath := upload.UploadParseAndDelKey(key)
	if len(filePath) != 2 {
		return int64(0), errors.New("这可能是个Bug!")
	}
	//这里有2个文件，分别是mmzz.h和mmzz.cpp.我们需要确认这两个文件都存在才算上传成功
	header := filePath[0]
	source := filePath[1]
	fileExist := oss.OssCheckFilesExist([]string{header, source})
	headerInfo := fileExist[header]
	sourceInfo := fileExist[source]
	if headerInfo == nil || sourceInfo == nil {
		return int64(0), errors.New("文件信息获取失败!")
	}
	//开启事务
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return int64(0), errors.New("开启session失败")
	}
	defer session.Rollback()
	//保存记录到数据库
	info := model.CodeAssessmentInfo{
		ID:         id,
		UploadTime: util.UTC_Time(),
		Header:     headerInfo.Key,
		Source:     sourceInfo.Key,
		Teacher:    teacherSelected,
		HeaderSize: headerInfo.Size,
		SourceSize: sourceInfo.Size,
	}
	indices := int64(0)
	if indices = dao.CodeAssessmentAdd(session, info); indices == 0 {
		return int64(0), errors.New("写入数据库失败!")
	}
	//提交事务
	if err := session.Commit(); err != nil {
		return int64(0), errors.New("commit事务失败")
	}
	//
	return indices, nil
}
