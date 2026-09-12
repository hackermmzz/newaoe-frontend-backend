package Email

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gopkg.in/gomail.v2"
)

// 消息队列里面的数据结构
type EmailMsg struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

var (
	EmailPushConsumer rocketmq.PushConsumer
	// Postfix SMTP Dialer
	// 不维持长连接，每次发送由 DialAndSend 建立连接。
	EmailDialer *gomail.Dialer
)

// 初始化邮件发送服务
func EmailSenderInit() {
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
	EmailPushConsumer = dao.NewMQPushConsumer(
		config.Conf.Email.EmailMQTopic,
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for _, msg := range msgs {
				var emailMsg EmailMsg
				if err := json.Unmarshal(msg.Body, &emailMsg); err != nil {
					util.Debug("邮件消息解析失败:", err)
					continue
				}
				if emailMsg.Email == "" ||
					emailMsg.Subject == "" ||
					emailMsg.Text == "" {
					util.Debug("邮件字段不完整")
					continue
				}
				if err := sendEmailMsgToServer(emailMsg); err != nil {
					util.Debug(
						"邮件提交到 Postfix 失败:",
						err,
						" 收件人:",
						emailMsg.Email,
					)
					return consumer.ConsumeRetryLater, nil
				}
			}
			return consumer.ConsumeSuccess, nil
		},
		groups...,
	)

	if err := EmailPushConsumer.Start(); err != nil {
		panic("EmailSenderInit:" + err.Error())
	}
}

// 将邮件提交给 Postfix
func sendEmailMsgToServer(emailMsg EmailMsg) error {

	m := WrapForTextEmail(emailMsg)
	m.SetHeader(
		"From",
		config.Conf.Email.SenderEmail,
	)
	err := EmailDialer.DialAndSend(m)
	if err != nil {
		util.Debug("提交邮件到 Postfix 失败:", err, " 收件人:", emailMsg.Email)
		return err
	}

	util.Debug("邮件已提交到 Postfix:", emailMsg.Email)

	return nil
}
