package Upload

import (
	"newaoe/dao"
	"newaoe/service/Code"
	"newaoe/util"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func ProcessUploadCode(ctx *gin.Context, userInfo map[string]interface{}, path []minio.ObjectInfo) {
	id, _ := userInfo["id"].(string)
	headerInfo := path[0]
	sourceInfo := path[1]
	descriptionInfo := path[2]
	description := dao.OssGetFileData(descriptionInfo.Key)
	//
	info := dao.CodeCommonInfo{
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
	session := dao.DB.NewSession()
	defer session.Rollback()
	err := session.Begin()
	if err != nil {
		util.DebugError(id, "代码上传失败", err.Error())
		util.ResponseNAK_MSG(ctx, "代码上传失败", "")
		return
	}
	//写入数据库
	if !dao.CodeCommonInfoAdd(session, info) {
		//
		util.DebugError(id, "代码上传失败", "中间过程有问题")
		util.ResponseNAK_MSG(ctx, "代码上传失败", "")
		return
	}
	//运行
	runinfo := dao.CodeRunInfo{
		ID:          id,
		SubmitTime:  util.UTC_Time(),
		Header:      headerInfo.Key,
		Source:      sourceInfo.Key,
		Class:       dao.Code_Common,
		Description: string(description),
		Status:      dao.ProcessDataMessageByStatus(dao.NewCodeRunStatusInfo()).Marshal(),
	}
	if !Code.RunUserCode(session, runinfo, true) {
		util.ResponseNAK_MSG(ctx, "代码上传失败", "")
		return
	}
	//提交事务
	if err = session.Commit(); err != nil {
		util.DebugError(id, "代码上传失败", err.Error())
		util.ResponseNAK_MSG(ctx, "代码上传失败", "")
		return
	}
	//
	util.ResponseACK_MSG(ctx, "代码上传成功", "")
	util.DebugSuccess(id, "代码上传成功!")
}

func CodeUploadConfirm(ctx *gin.Context, userInfo map[string]interface{}, path []string) {
	//这里有三个文件，分别是mmzz.h和mmzz.cpp和mmzz.txt
	//我们需要确认这几个文件都存在才算上传成功
	header := path[0]
	source := path[1]
	description := path[2]
	fileExists := dao.OssCheckFilesExist([]string{header, source, description})
	headerInfo := fileExists[header]
	sourceInfo := fileExists[source]
	descriptionInfo := fileExists[description]
	if headerInfo == nil || sourceInfo == nil || descriptionInfo == nil {
		util.ResponseNAK_MSG(ctx, "上传文件确认失败", "")
		return
	}
	//处理上传的代码文件
	ProcessUploadCode(ctx, userInfo, []minio.ObjectInfo{*headerInfo, *sourceInfo, *descriptionInfo})
}
