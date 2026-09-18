package controller

import (
	"newaoe/Src/codeCommonSubmit/dao"
	"newaoe/Src/codeCommonSubmit/model"
	"newaoe/Src/common/upload"
	database "newaoe/Src/databse"
	"newaoe/Src/oss"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func CodeCommonSubmitACK(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//获取对应的文件路径
	data_bytes, _ := ctx.GetRawData()
	filePath := upload.UploadParseAndDelKey(data_bytes, "key", id)
	if len(filePath) != 3 {
		util.DebugError("这可能是一个Bug!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//获取信息
	fileInfos := oss.OssCheckFilesExist(filePath)
	if len(fileInfos) != len(filePath) {
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//
	headerInfo, ok0 := fileInfos[filePath[0]]
	sourceInfo, ok1 := fileInfos[filePath[1]]
	descriptionInfo, ok2 := fileInfos[filePath[2]]
	if !ok1 || !ok2 || !ok0 || descriptionInfo == nil || sourceInfo == nil || headerInfo == nil {
		util.DebugError("为什么获取不到对应的信息!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
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
		util.DebugError(id, "代码上传失败", err.Error())
		util.ResponseNAK_MSG(ctx, "代码上传失败", "")
		return
	}
	defer session.Rollback()
	//写入数据库
	indices := int64(0)
	if indices = dao.CodeCommonInfoAdd(session, info); indices == 0 {
		//
		util.DebugError(id, "代码上传失败中间过程有问题")
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
	util.ResponseACK_MSG(ctx, "代码上传成功", map[string]interface{}{"indices": indices})
	util.DebugSuccess(id, "代码上传成功!")
}
