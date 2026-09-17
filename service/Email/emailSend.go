package Email

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"

	"github.com/apache/rocketmq-client-go/v2/primitive"
)

// 发送邮件(推送到消息队列)
func SendEmail(email EmailMsg) {
	//写入数据库
	session := dao.DB.NewSession()
	if err := session.Begin(); err != nil {
		util.DebugError("数据库异常!")
		return
	}
	defer session.Close()
	indices := dao.EmailInfoInsert(session, dao.EmailInfo{
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
	dao.RocketMQProducer.SendAsync(context.Background(),
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
