package Upload

import (
	"fmt"
	"newaoe/dao"
	"newaoe/service/Code"
	data "newaoe/service/Data"
	"newaoe/util"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func AssessmentSubmissionUploadConfirm(ctx *gin.Context, userInfo map[string]string, path []string) {
	//这里有2个文件，分别是mmzz.h和mmzz.cpp
	//我们需要确认这两个文件都存在才算上传成功
	header := path[0]
	source := path[1]
	fileExist := dao.OssCheckFilesExist([]string{header, source})
	headerInfo := fileExist[header]
	sourceInfo := fileExist[source]
	if headerInfo == nil || sourceInfo == nil {
		util.ResponseNAK_MSG(ctx, "上传文件确认失败", "")
		return
	}
	//处理上传的代码文件
	AssessmentSubmissionUploadProcess(ctx, userInfo, []minio.ObjectInfo{*headerInfo, *sourceInfo})
}

func AssessmentSubmissionUploadProcess(ctx *gin.Context, userInfo map[string]string, file []minio.ObjectInfo) {
	//
	id := userInfo["id"]
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
	if !data.TeacherExist(teacher) {
		util.DebugError("AssessmentSubmissionUploadProcess:这里有bug")
		util.ResponseNAK_MSG(ctx, "你小子怎么绕过来的", "")
		return
	}
	//开启事务
	session := dao.DB.NewSession()
	defer session.Rollback()
	err := session.Begin()
	if err != nil {
		util.DebugError("AssessmentSubmissionUploadProcess:服务器异常!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", "")
		return
	}
	//保存记录到数据库
	header := file[0]
	source := file[1]
	if !dao.CodeAssessmentAdd(session, dao.CodeAssessmentInfo{
		ID:         id,
		UploadTime: util.UTC_Time(),
		Header:     header.Key,
		Source:     source.Key,
		Teacher:    teacher,
		HeaderSize: header.Size,
		SourceSize: source.Size,
	}) {
		util.DebugError("AssessmentSubmissionUploadProcess:记录数据库出错!")
		util.ResponseNAK_MSG(ctx, "上传出错!", "")
		return
	}
	//运行代码
	info := dao.CodeRunInfo{
		ID:          id,
		SubmitTime:  util.UTC_Time(),
		Header:      header.Key,
		Source:      source.Key,
		Class:       dao.Code_Common,
		Description: fmt.Sprintf("this is the final code you submit (teacher:%v)", teacher),
		Status:      dao.ProcessDataMessageByStatus(dao.NewCodeRunStatusInfo()).Marshal(),
	}
	if !Code.RunUserCode(session, info) {
		util.ResponseNAK_MSG(ctx, "上传失败!", "")
		return
	}
	//提交事务
	if err := session.Commit(); err != nil {
		util.ResponseNAK_MSG(ctx, "服务器异常!", "")
		return
	}
	//
	util.ResponseACK_MSG(ctx, "上传成功!", "")
	util.DebugSuccess(id, "上传考核代码成功!")
	//
}
