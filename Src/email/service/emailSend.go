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

	"github.com/apache/rocketmq-clients/golang/v5"
	rocketmq "github.com/apache/rocketmq-clients/golang/v5"
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
	EmailPushConsumer golang.PushConsumer
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
	EmailPushConsumer = mq.NewMQPushConsumer(
		config.Conf.Email.EmailMQTopic,
		func(msg *rocketmq.MessageView) rocketmq.ConsumerResult {
			return consumeEmailToEmailService(msg.GetBody())
		},
		"EmailMQ_group",
	)
	if err := EmailPushConsumer.Start(); err != nil {
		panic("EmailSenderInit:" + err.Error())
	}
	//
}

// 发送邮件(推送到消息队列)
func SendEmail(email EmailMsg) {
	//写入数据库
	session := database.NewSession()
	if err := session.Begin(); err != nil {
		util.DebugError("SendEmail:事务开始失败!")
		return
	}
	defer session.Close()
	indices, err := dao.EmailInfoInsert(session, model.EmailInfo{
		Receiver:   email.Email,
		CreateTime: util.UTC_Time(),
		Class:      email.Type,
		Data:       util.TruncateString(email.Text, 16*1024), //截断为16kb
		Send:       false,
	})
	if err != nil {
		util.DebugError("EmailInfoInsert:", err)
	}
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
	msg := &rocketmq.Message{
		Topic: config.Conf.Email.EmailMQTopic,
		Body:  data_byte,
	}
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
func sendEmailWithRetry(msg *rocketmq.Message, limit int) {
	if limit <= 0 {
		util.DebugError("sendEmailWithRetry达到最大发送次数!", (string)(msg.Body))
		return
	}
	//放入队列
	mq.RocketMQProducer.SendAsync(
		context.Background(),
		msg,
		func(ctx context.Context, _ []*rocketmq.SendReceipt, err error) {
			if err != nil {
				util.DebugError("sendEmailWithRetry:", err)
			}
		},
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

func consumeEmailToEmailService(data []byte) rocketmq.ConsumerResult {
	var emailMsg EmailMsgInQueue
	//解析
	if err := json.Unmarshal(data, &emailMsg); err != nil {
		util.DebugError("邮件消息解析失败:", err)
		return rocketmq.SUCCESS
	}
	//检查字段
	if emailMsg.Indices == 0 || emailMsg.Email == "" || emailMsg.Subject == "" || emailMsg.Text == "" {
		util.DebugError("邮件字段不完整")
		return rocketmq.SUCCESS
	}
	//判断是否已经发送过了
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		util.DebugError("数据库异常!")
		return rocketmq.FAILURE
	}
	defer session.Rollback()
	info, err := dao.EmailInfoGetForUpdate(session, emailMsg.Indices)
	if err != nil {
		util.DebugError("EmailInfoGetForUpdate:", err)
	}
	if info == nil {
		util.DebugError("这是一个bug!按道理不应该为nil")
		return rocketmq.SUCCESS
	}
	if info.Send {
		return rocketmq.SUCCESS
	}
	//发送邮件
	if err := sendEmailMsgToServer(emailMsg.EmailMsg); err != nil {
		util.DebugError("邮件提交到 Postfix 失败:", err, " 收件人:", emailMsg.Email)
		return rocketmq.FAILURE
	}
	//标记已经发送过了
	ok, _, err := dao.EmailInfoUpdateSendStatus(session, info.Indices, true)
	if err != nil {
		util.DebugError("EmailInfoUpdateSendStatus:", err)
	}
	if !ok {
		util.DebugError("send字段修改失败!")
		return rocketmq.FAILURE
	}
	if err := session.Commit(); err != nil {
		util.DebugError("邮件send标记位修改失败!", err)
		return rocketmq.FAILURE
	}
	return rocketmq.SUCCESS
}
