package Home

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/service/Email"
	"newaoe/util"
	"path"
	"strconv"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
)

var FeedBackEmailPushConsumer rocketmq.PushConsumer

func FeedbackInit() {
	//一个组足矣
	var groups []string
	for i := 0; i < 1; i++ {
		groups = append(groups, "FeedbackEmailMQ_group"+strconv.Itoa(i))
	}
	//
	FeedBackEmailPushConsumer = dao.NewMQPushConsumer(
		config.Conf.Feedback.FeedbackNeedSendToEmailTopic,
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				link := string(msg.Body)
				//送到邮箱队列
				email := Email.EmailMsg{
					Email:   config.Conf.Feedback.FeedbackSendToEmail,
					Text:    link,
					Subject: "NewAOE网站反馈",
					Type:    Email.EmailMsgType_TEXT,
				}
				data_byte, _ := json.Marshal(email)
				//同步发
				_, err := dao.RocketMQProducer.SendSync(
					context.Background(),
					primitive.NewMessage(
						config.Conf.Email.EmailMQTopic,
						data_byte,
					),
				)
				if err == nil {
					util.DebugSuccess("反馈邮件成功推送到邮箱!")
				} else {
					util.DebugError("FeedBackEmailPushConsumer send to emailqueue fail!:", err)
					return consumer.ConsumeRetryLater, nil
				}
			}
			return consumer.ConsumeSuccess, nil
		},
		groups...,
	)
	FeedBackEmailPushConsumer.Start()
}

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
	session := dao.DB.NewSession()
	defer session.Close()
	feedbackinfo := dao.FeedbackGetByIndices(session, postInfo.Indices)
	if feedbackinfo == nil {
		util.DebugError("数据库获取反馈记录失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//上传html数据
	htmlpath := path.Join(feedbackinfo.BaseFolder, "feedback.html")
	if !dao.OssUploadFileData(htmlpath, []byte(postInfo.Text), "text/html; charset=utf-8") {
		util.DebugError("上传html失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//获取链接
	url := dao.OssGetDownloadFileUrl(htmlpath, time.Duration(7*24)*time.Hour, false)
	if url == "" {
		util.DebugError("获取html下载链接失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//发送给我的邮件
	data := fmt.Sprintf("反馈访问链接: %s", url)
	msg := primitive.NewMessage(config.Conf.Feedback.FeedbackNeedSendToEmailTopic, []byte(data))
	sendFeedbackToEmail(msg, 3) //最多尝试三次
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
	imagesurl := dao.GetUploadFileUrls(imagespath, util.NewArray(len(data.Images), expire_time))
	filesurl := dao.GetUploadFileUrls(filespath, util.NewArray(len(data.Files), expire_time))
	videosurl := dao.GetUploadFileUrls(videospath, util.NewArray(len(data.Videos), expire_time))
	//生成下载链接
	downloadExpireTime := time.Duration(7*24) * time.Hour
	imagesDownloadurls := dao.OssGetDownloadFileUrls(imagespath, downloadExpireTime, false)
	filesDownloadUrls := dao.OssGetDownloadFileUrls(filespath, downloadExpireTime, true)
	videosDownloadUrls := dao.OssGetDownloadFileUrls(videospath, downloadExpireTime, false)
	//
	if len(imagesurl) != len(data.Images) || len(filesurl) != len(data.Files) ||
		len(imagesDownloadurls) != len(data.Images) || len(filesDownloadUrls) != len(data.Files) ||
		len(videosDownloadUrls) != len(data.Videos) {
		util.DebugError("上传/下载链接生成失败!")
		util.ResponseNAK_MSG(ctx, "服务器异常!", nil)
		return
	}
	//插入记录到数据库
	session := dao.DB.NewSession()
	defer session.Close()
	info := dao.FeedbackInfo{
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

func sendFeedbackToEmail(msg *primitive.Message, limit int) {
	if limit <= 0 {
		util.DebugError("sendFeedbackToEmail超出重试次数!")
		return
	}
	//
	dao.RocketMQProducer.SendAsync(context.Background(),
		func(ctx context.Context, result *primitive.SendResult, err error) {
			// 回调：发送完才进来
			if err != nil {
				util.DebugError("反馈推送失败!", err)
				sendFeedbackToEmail(msg, limit-1)
			} else {
				util.DebugSuccess("反馈推送成功!")
			}
		},
		msg,
	)
}
