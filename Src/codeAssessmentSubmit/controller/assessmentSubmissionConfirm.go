package controller

import (
	"newaoe/Src/codeAssessmentSubmit/dao"
	"newaoe/Src/codeAssessmentSubmit/model"
	"newaoe/Src/common/upload"
	database "newaoe/Src/databse"
	"newaoe/Src/oss"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

func AssessmentSubmissionSubmitACK(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//获取对应的文件路径
	data_bytes, _ := ctx.GetRawData()
	filePath := upload.UploadParseAndDelKey(data_bytes, "key", id)
	if len(filePath) != 2 {
		util.DebugError("这可能是一个Bug!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//这里有2个文件，分别是mmzz.h和mmzz.cpp.我们需要确认这两个文件都存在才算上传成功
	header := filePath[0]
	source := filePath[1]
	fileExist := oss.OssCheckFilesExist([]string{header, source})
	headerInfo := fileExist[header]
	sourceInfo := fileExist[source]
	if headerInfo == nil || sourceInfo == nil {
		util.ResponseNAK_MSG(ctx, "上传文件确认失败", "")
		return
	}
	//获取任课老师
	type Body struct {
		Teacher string `json:"teacher"`
	}
	var body Body
	if !util.JsonCtx(ctx, &body) {
		util.DebugError("AssessmentSubmissionUploadProcess:有人伪造请求!")
		util.ResponseNAK_MSG(ctx, "你小子怎么绕过来的", "")
		return
	}
	teacher := body.Teacher
	if !dao.TeacherExist(teacher) {
		util.DebugError("AssessmentSubmissionUploadProcess:这里有bug")
		util.ResponseNAK_MSG(ctx, "你小子怎么绕过来的", "")
		return
	}
	//开启事务
	session := database.NewSession()
	err := session.Begin()
	if err != nil {
		util.DebugError("AssessmentSubmissionUploadProcess:服务器异常!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", "")
		return
	}
	defer session.Rollback()
	//保存记录到数据库
	info := model.CodeAssessmentInfo{
		ID:         id,
		UploadTime: util.UTC_Time(),
		Header:     headerInfo.Key,
		Source:     sourceInfo.Key,
		Teacher:    teacher,
		HeaderSize: headerInfo.Size,
		SourceSize: sourceInfo.Size,
	}
	indices := int64(0)
	if indices = dao.CodeAssessmentAdd(session, info); indices == 0 {
		util.DebugError("AssessmentSubmissionUploadProcess:记录数据库出错!")
		util.ResponseNAK_MSG(ctx, "上传出错!", "")
		return
	}
	//提交事务
	if err := session.Commit(); err != nil {
		util.ResponseNAK_MSG(ctx, "服务器异常!", "")
		return
	}
	//
	util.ResponseACK_MSG(ctx, "上传成功!", map[string]interface{}{"indices": indices})
	util.DebugSuccess(id, "上传考核代码成功!")
}
