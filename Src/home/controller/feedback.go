package controller

import (
	"fmt"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/email/service"
	"newaoe/Src/home/dao"
	"newaoe/Src/home/model"
	"newaoe/Src/oss"
	"newaoe/Src/util"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func StudentFeedbackUpload(ctx *gin.Context) {
	//解析数据
	type PostDataInfo struct {
		Indices int    `json:"indices"`
		Text    string `json:"text"`
	}
	var postInfo PostDataInfo
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
	//获取反馈的基础路径
	session := database.NewSession()
	defer session.Close()
	feedbackinfo := dao.FeedbackGetByIndices(session, postInfo.Indices)
	if feedbackinfo == nil {
		util.DebugError("数据库获取反馈记录失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//上传html数据
	htmlpath := path.Join(feedbackinfo.BaseFolder, "feedback.html")
	if !oss.OssUploadFileData(htmlpath, []byte(postInfo.Text), "text/html; charset=utf-8") {
		util.DebugError("上传html失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//获取链接
	url := oss.OssGetDownloadFileUrl(htmlpath, time.Duration(7*24)*time.Hour, false)
	if url == "" {
		util.DebugError("获取html下载链接失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//送到邮箱队列
	link := fmt.Sprintf("反馈访问链接: %s", url)
	email := service.EmailMsg{
		Email:   config.Conf.Feedback.FeedbackSendToEmail,
		Text:    link,
		Subject: "NewAOE网站反馈",
		Type:    service.EmailMsgType_TEXT,
	}
	service.SendEmail(email)
	//
	util.ResponseACK_MSG(ctx, "发送成功!", nil)
}

func StudentFeedbckUploadAttachment(ctx *gin.Context) {
	//获取用户id
	userInfo := util.GetCtxTookenInfo(ctx)
	if userInfo == nil {
		util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
		return
	}
	//解析数据
	type ReceivePostDataInfo struct {
		Images []string `json:"images"`
		Videos []string `json:"videos"`
		Files  []string `json:"files"`
	}
	var data ReceivePostDataInfo
	if !util.JsonCtx(ctx, &data) {
		util.ResponseNAK_MSG(ctx, "格式错误!", nil)
		return
	}
	if len(data.Images) > 20 || len(data.Files) >= 20 || len(data.Videos) >= 5 {
		util.ResponseNAK_MSG(ctx, "传输文件过多!", nil)
		return
	}
	//生成文件名
	id := userInfo["id"].(string)
	dirName := id + "_" + util.UTC_Time().Format("20060102_150405")
	dir := path.Join(config.Conf.OSS.PublicBaseFolder, config.Conf.Feedback.FeedbackFolder, dirName)
	imagespath := make([]string, len(data.Images))
	filespath := make([]string, len(data.Files))
	videospath := make([]string, len(data.Videos))
	for i := range imagespath {
		imagespath[i] = fmt.Sprintf("%v/image_%d_%s", dir, i, data.Images[i])
	}
	for i := range filespath {
		filespath[i] = fmt.Sprintf("%v/file_%d_%s", dir, i, data.Files[i])
	}
	for i := range videospath {
		videospath[i] = fmt.Sprintf("%v/video_%d_%s", dir, i, data.Videos[i])
	}
	//生成上传链接
	expire_time := time.Duration(60) * time.Minute
	imagesurl := oss.GetUploadFileUrls(imagespath, util.NewArray(len(data.Images), expire_time))
	filesurl := oss.GetUploadFileUrls(filespath, util.NewArray(len(data.Files), expire_time))
	videosurl := oss.GetUploadFileUrls(videospath, util.NewArray(len(data.Videos), expire_time))
	//生成下载链接
	downloadExpireTime := time.Duration(7*24) * time.Hour
	imagesDownloadurls := oss.OssGetDownloadFileUrls(imagespath, downloadExpireTime, false)
	filesDownloadUrls := oss.OssGetDownloadFileUrls(filespath, downloadExpireTime, true)
	videosDownloadUrls := oss.OssGetDownloadFileUrls(videospath, downloadExpireTime, false)
	//
	if len(imagesurl) != len(data.Images) || len(filesurl) != len(data.Files) ||
		len(imagesDownloadurls) != len(data.Images) || len(filesDownloadUrls) != len(data.Files) ||
		len(videosDownloadUrls) != len(data.Videos) {
		util.DebugError("上传/下载链接生成失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//插入记录到数据库
	session := database.NewSession()
	defer session.Close()
	info := model.FeedbackInfo{
		SubmitTime: util.UTC_Time(),
		BaseFolder: dir,
	}
	indices := dao.FeedbackInsert(session, info)
	if indices == 0 {
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	if err := session.Commit(); err != nil {
		util.DebugError("StudentFeedbckUploadAttachment:", err)
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//返回数据
	type ReplyInfo struct {
		Indices           int      `json:"indices"`
		UploadImageUrls   []string `json:"uploadimageurls"`
		UploadFileUrls    []string `json:"uploadfileurls"`
		UploadVideoUrls   []string `json:"uploadvideourls"`
		DownloadImageUrls []string `json:"downloadimageurls"`
		DownloadFileUrls  []string `json:"downloadfileurls"`
		DownloadVideoUrls []string `json:"downloadvideourls"`
	}
	replyInfo := ReplyInfo{
		Indices:           indices,
		UploadImageUrls:   imagesurl,
		UploadFileUrls:    filesurl,
		UploadVideoUrls:   videosurl,
		DownloadImageUrls: imagesDownloadurls,
		DownloadFileUrls:  filesDownloadUrls,
		DownloadVideoUrls: videosDownloadUrls,
	}
	util.ResponseACK_MSG(ctx, "成功", replyInfo)
}
