package Email

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gopkg.in/gomail.v2"
)

func encodeChinese(name string) string {
	// 用 Base64 编码中文（邮件头推荐用 Base64 编码）
	b64Str := base64.StdEncoding.EncodeToString([]byte(name))
	return "=?UTF-8?B?" + b64Str + "?="
}

// 发送简单文本邮件(每次发送都重新建立连接)
func SendTextEmail(recvEmail string, subject string, text string) {
	email := EmailMsg{
		Email:   recvEmail,
		Text:    text,
		Subject: subject,
	}
	//放入队列
	data_byte, _ := json.Marshal(&email)
	dao.RocketMQProducer.SendAsync(context.Background(),
		func(ctx context.Context, result *primitive.SendResult, err error) {
			// 回调：发送完才进来
			if err != nil {
				util.Debug("SendTextEmail:", err)
				return
			}
		},
		primitive.NewMessage(config.Conf.Email.EmailMQTopic, data_byte),
	)
}

// 包装email为可发送对象
func WrapForTextEmail(email EmailMsg) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", encodeChinese("提瓦特须弥智慧之神纳西妲")+" <"+config.Conf.Email.SenderEmail+">")
	msg.SetHeader("To", email.Email)
	msg.SetHeader("Subject", email.Subject)
	msg.SetBody("text/plain", email.Text)
	return msg
}
