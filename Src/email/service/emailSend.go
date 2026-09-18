package service

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/email/dao"
	"newaoe/Src/email/model"
	"newaoe/Src/mq"
	"newaoe/Src/util"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gopkg.in/gomail.v2"
)

var (
	EmailMsgType_TEXT = 0
	EmailMsgType_XML  = 1
)

// 用户传入得Email类型
type EmailMsg struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	Type    int    `json:"type"`
}

// 消息队列的Email类型
type EmailMsgInQueue struct {
	Indices int `json:"indices"`
	EmailMsg
}

var (
	EmailPushConsumer rocketmq.PushConsumer
	// Postfix SMTP Dialer
	// 不维持长连接，每次发送由 DialAndSend 建立连接。
	EmailDialer *gomail.Dialer
)

// 初始化邮件发送服务
func EmailSendServiceInit() {
	//
	EmailDialer = gomail.NewDialer(
		config.Conf.Email.EmailServe,
		config.Conf.Email.EmailServePort,
		"",
		"",
	)
	EmailDialer.SSL = false
	// backend -> postfix 是 Docker 内网
	// Postfix 使用自签名证书，因此关闭该段证书校验
	EmailDialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	var groups []string
	for i := 0; i < config.Conf.Email.EmailMQGroupCount; i++ {
		groups = append(groups, "EmailMQ_group"+strconv.Itoa(i))
	}
	EmailPushConsumer = mq.NewMQPushConsumer(
		config.Conf.Email.EmailMQTopic,
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				if ret := consumeEmailToEmailService(msg.Body); ret != consumer.ConsumeSuccess {
					return ret, nil
				}
			}
			return consumer.ConsumeSuccess, nil
		},
		groups...,
	)

	if err := EmailPushConsumer.Start(); err != nil {
		panic("EmailSenderInit:" + err.Error())
	}
	//
}

// 发送邮件(推送到消息队列)
func SendEmail(email EmailMsg) {
	//写入数据库
	session := database.DB.NewSession()
	if err := session.Begin(); err != nil {
		util.DebugError("数据库异常!")
		return
	}
	defer session.Close()
	indices := dao.EmailInfoInsert(session, model.EmailInfo{
		Receiver:   email.Email,
		CreateTime: util.UTC_Time(),
		Class:      email.Type,
		Data:       util.TruncateString(email.Text, 16*1024), //截断为16kb
		Send:       false,
	})
	if indices == 0 {
		util.DebugError("SendEmail插入记录失败!")
		return
	}
	if err := session.Commit(); err != nil {
		util.DebugError("SendEmail事务提交失败!")
		return
	}
	//构造数据
	emailInQUeue := EmailMsgInQueue{
		EmailMsg: email,
		Indices:  indices,
	}
	data_byte, _ := json.Marshal(&emailInQUeue)
	msg := primitive.NewMessage(config.Conf.Email.EmailMQTopic, data_byte)
	//默认重试三次
	go sendEmailWithRetry(msg, 3)
}

// 编码为中文
func encodeChinese(name string) string {
	// 用 Base64 编码中文（邮件头推荐用 Base64 编码）
	b64Str := base64.StdEncoding.EncodeToString([]byte(name))
	return "=?UTF-8?B?" + b64Str + "?="
}

// 重试发送
func sendEmailWithRetry(msg *primitive.Message, limit int) {
	if limit <= 0 {
		util.DebugError("sendEmailWithRetry达到最大发送次数!", (string)(msg.Body))
		return
	}
	//放入队列
	mq.RocketMQProducer.SendAsync(context.Background(),
		func(ctx context.Context, result *primitive.SendResult, err error) {
			// 回调：发送完才进来
			if err != nil {
				util.DebugError("sendEmailWithRetry:", err)
				return
			}
		},
		msg,
	)
}

// 将邮件提交给 Postfix
func sendEmailMsgToServer(emailMsg EmailMsg) error {
	var msg *gomail.Message
	//分类
	switch emailMsg.Type {
	case EmailMsgType_TEXT:
		msg = wrapForTextEmail(emailMsg)
	case EmailMsgType_XML:
		msg = wrapForHTMLEmail(emailMsg)
	}
	//
	msg.SetHeader("From", config.Conf.Email.SenderEmail)
	err := EmailDialer.DialAndSend(msg)
	if err != nil {
		util.DebugError("提交邮件到 Postfix 失败:", err, " 收件人:", emailMsg.Email)
		return err
	}

	util.DebugSuccess("邮件已提交到 Postfix:", emailMsg.Email)

	return nil
}

// 包装email为可发送对象
func wrapForTextEmail(email EmailMsg) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", encodeChinese("提瓦特须弥智慧之神纳西妲")+" <"+config.Conf.Email.SenderEmail+">")
	msg.SetHeader("To", email.Email)
	msg.SetHeader("Subject", email.Subject)
	msg.SetBody("text/plain", email.Text)
	return msg
}

// WrapForHTMLEmail 网络外链图片版本
func wrapForHTMLEmail(email EmailMsg) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", encodeChinese("提瓦特须弥智慧之神纳西妲")+" <"+config.Conf.Email.SenderEmail+">")
	msg.SetHeader("To", email.Email)
	msg.SetHeader("Subject", email.Subject)

	// multipart/alternative：纯文本降级 + HTML富文本
	msg.SetBody("text/plain", email.Text)
	msg.AddAlternative("text/html", email.Text)
	msg.AddAlternative("text/html", email.Text, gomail.SetPartEncoding(gomail.Base64))
	return msg
}

func consumeEmailToEmailService(data []byte) consumer.ConsumeResult {
	var emailMsg EmailMsgInQueue
	//解析
	if err := json.Unmarshal(data, &emailMsg); err != nil {
		util.DebugError("邮件消息解析失败:", err)
		return consumer.ConsumeSuccess
	}
	//检查字段
	if emailMsg.Indices == 0 || emailMsg.Email == "" || emailMsg.Subject == "" || emailMsg.Text == "" {
		util.DebugError("邮件字段不完整")
		return consumer.ConsumeSuccess
	}
	//判断是否已经发送过了
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		util.DebugError("数据库异常!")
		return consumer.ConsumeRetryLater
	}
	defer session.Rollback()
	info := dao.EmailInfoGetForUpdate(session, emailMsg.Indices)
	if info == nil {
		util.DebugError("这是一个bug!按道理不应该为nil")
		return consumer.ConsumeSuccess
	}
	if info.Send {
		return consumer.ConsumeSuccess
	}
	//发送邮件
	if err := sendEmailMsgToServer(emailMsg.EmailMsg); err != nil {
		util.DebugError("邮件提交到 Postfix 失败:", err, " 收件人:", emailMsg.Email)
		return consumer.ConsumeRetryLater
	}
	//标记已经发送过了
	ok, _ := dao.EmailInfoUpdateSendStatus(session, info.Indices, true)
	if !ok {
		util.DebugError("send字段修改失败!")
		return consumer.ConsumeRetryLater
	}
	if err := session.Commit(); err != nil {
		util.DebugError("邮件send标记位修改失败!", err)
		return consumer.ConsumeRetryLater
	}
	return consumer.ConsumeSuccess
}
