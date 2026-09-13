package Home

import (
	"context"
	"encoding/json"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/service/Email"
	"newaoe/service/Upload"
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
		config.Conf.Other.FeedbackNeedSendToEmailTopic,
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				//解析文件地址
				var path string
				if err := json.Unmarshal(msg.Body, &path); err != nil {
					util.Debug("反馈信息解析失败:", err)
					continue
				}
				//获取文件链接
				expireDuration := time.Duration(config.Conf.Other.PrivateFileDownloadUrlExpireTime) * time.Second
				url := dao.OssGetDownloadFileUrl(path, expireDuration, false)
				if url == "" {
					util.Debug("FeedBackEmailPushConsumer get download link fail!")
					return consumer.ConsumeRetryLater, nil
				}
				util.Debug(config.Conf.Other.FeedbackSendToEmail)
				//送到邮箱队列
				email := Email.EmailMsg{
					Email:   config.Conf.Other.FeedbackSendToEmail,
					Text:    "有学生发送反馈消息,链接: " + url,
					Subject: "NewAOE网站反馈",
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
					util.Debug("邮件成功推送到邮箱!")
				} else {
					util.Debug("FeedBackEmailPushConsumer send to emailqueue fail!:", err)
					return consumer.ConsumeRetryLater, nil
				}
			}
			return consumer.ConsumeSuccess, nil
		},
		groups...,
	)
	FeedBackEmailPushConsumer.Start()
}

func StudentFeedback(ctx *gin.Context) {
	//获取用户id
	userInfo := util.GetCtxTookenInfo(ctx)
	if userInfo == nil {
		util.ResponseNAK_MSG(ctx, "cookie过期或者错误!", "")
		return
	}
	id := userInfo["id"]
	//指定文件名称
	filename := util.UTC_Time().Format("20060102_150405") + "_" + id + ".json"
	targetfile := path.Join(config.Conf.OSS.PublicBaseFolder, "feedback", filename)
	//走upload的路线,但是不会confirm(url过期时间和普通代码提交时间一样)
	duration := time.Duration(24*7) * time.Hour //一周就行了
	if Upload.UploadFile(ctx, []string{targetfile}, []time.Duration{duration}) {
		//发送给我的邮件
		data_byte, _ := json.Marshal(targetfile)
		dao.RocketMQProducer.SendAsync(context.Background(),
			func(ctx context.Context, result *primitive.SendResult, err error) {
				// 回调：发送完才进来
				if err != nil {
					util.Debug("反馈推送失败!", err)
					return
				} else {
					util.Debug("反馈推送成功!")
				}
			},
			primitive.NewMessage(config.Conf.Other.FeedbackNeedSendToEmailTopic, data_byte),
		)
	}
}
