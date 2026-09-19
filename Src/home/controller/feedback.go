package controller

import (
	"fmt"
	"newaoe/Src/config"
	emailService "newaoe/Src/email/service"
	"newaoe/Src/home/service"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

type postDataInfo struct {
	Indices int    `json:"indices"`
	Text    string `json:"text"`
}

func StudentFeedbackUpload(ctx *gin.Context) {
	//解析数据
	var postInfo postDataInfo
	if !util.JsonCtx(ctx, &postInfo) {
		util.DebugError("格式解析错误!")
		util.ResponseNAK_MSG(ctx, "格式错误!", nil)
		return
	}
	//如果html过长(大于1MB)，返回报错
	if len(postInfo.Text) > 1024*1024 {
		util.DebugError("反馈太长了!")
		util.ResponseNAK_MSG(ctx, "反馈过长!", nil)
		return
	}
	//获取反馈的下载链接
	url, err := service.FeedbackGetDownloadLink(postInfo.Indices, postInfo.Text)
	if err != nil {
		util.DebugError("StudentFeedbackUpload", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//送到邮箱队列
	link := fmt.Sprintf("反馈访问链接: %s", url)
	email := emailService.EmailMsg{
		Email:   config.Conf.Feedback.FeedbackSendToEmail,
		Text:    link,
		Subject: "NewAOE网站反馈",
		Type:    emailService.EmailMsgType_TEXT,
	}
	emailService.SendEmail(email)
	//
	util.ResponseACK_MSG(ctx, "发送成功!", nil)
}

type receivePostDataInfo struct {
	Images []string `json:"images"`
	Videos []string `json:"videos"`
	Files  []string `json:"files"`
}

func StudentFeedbckUploadAttachment(ctx *gin.Context) {
	//获取用户id
	userInfo := util.GetCtxTookenInfo(ctx)
	if userInfo == nil {
		util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
		return
	}
	//解析数据
	var data receivePostDataInfo
	if !util.JsonCtx(ctx, &data) {
		util.ResponseNAK_MSG(ctx, "格式错误!", nil)
		return
	}
	if len(data.Images) > 20 || len(data.Files) >= 20 || len(data.Videos) >= 5 {
		util.ResponseNAK_MSG(ctx, "传输文件过多!", nil)
		return
	}

	//返回数据
	replyInfo, err := service.FeedbackGetUploadUrls(userInfo["id"].(string), data.Images, data.Files, data.Videos)
	if err != nil {
		util.DebugError("StudentFeedbckUploadAttachment", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//
	util.ResponseACK_MSG(ctx, "成功", replyInfo)
}
