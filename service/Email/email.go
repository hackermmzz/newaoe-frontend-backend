package Email

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"
	"strconv"
	"errors"
	"io"
	"net"
	"strings"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"gopkg.in/gomail.v2"
)

// 消息队列里面数据结构体
type EmailMsg struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

// 待发送队列
var (
	SenderHanlder     gomail.SendCloser
	EmailPushConsumer rocketmq.PushConsumer
)


func EmailSenderInit() {
	connectToEmailServer()
	//初始化Consumer
	var err error
	var groups []string
	for i := 0; i < config.Conf.Email.EmailMQGroupCount; i++ {
		groups = append(groups, "EmailMQ_group"+strconv.Itoa(i))
	}

	EmailPushConsumer = dao.NewMQPushConsumer(config.Conf.Email.EmailMQTopic, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			var emailMsg EmailMsg
			//安全解析 JSON（必须判断错误）
			err := json.Unmarshal(msg.Body, &emailMsg)
			if err != nil {
				util.Debug("邮件消息解析失败:", err, " 原始内容:", string(msg.Body))
				continue
			}
			//非空判断（防止空指针）
			if emailMsg.Email == "" || emailMsg.Subject == "" || emailMsg.Text == "" {
				util.Debug("邮件字段不完整")
				continue
			}
			//连续重发三次
			for i:=0;i<3;i=i+1{
				err=sendEmailMsgToServer(emailMsg)
				if(err!=nil&&emailSevrerNeedReconnect(err)){
					util.Debug("Emial connection closed and now I reopen it!");
					connectToEmailServer()
				}else{
					break
				}
			}
		}

		// 5. 必须返回成功
		return consumer.ConsumeSuccess, nil
	}, groups...)
	err = EmailPushConsumer.Start()
	if err != nil {
		panic("EmailSenderInit:" + err.Error())
	}
	//运行发送邮件服务
}

func connectToEmailServer(){
	d := gomail.NewDialer(
		config.Conf.Email.EmailServe,
		config.Conf.Email.EmailServePort,
		config.Conf.Email.SenderEmail,
		config.Conf.Email.SenderAuthCode,
	)
	d.SSL = true

	// ====================== 关键修复：关闭 TLS 证书验证 ======================
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: true, // 这一行解决 Docker x509 报错
	}
	//
	var err error
	// 建立持久连接(尝试10次)
	for i := 0; i < 10; i++ {
		SenderHanlder, err = d.Dial()
		if err == nil {
			util.Debug("邮箱连接成功")
			return 
		}
	}
	panic("邮箱连接失败:" + err.Error())
}


func sendEmailMsgToServer(emailMsg EmailMsg) error{

	//发送邮件
	err:= gomail.Send(SenderHanlder, WrapForTextEmail(emailMsg))
	
	if err != nil {
		util.Debug("发送邮件出错:", err)
		//
		return err
	} 
	util.Debug("发送邮件成功:", emailMsg.Email)
	return nil
}

func emailSevrerNeedReconnect(err error) bool {
	if err == nil {
		return false
	}

	// EOF 通常意味着 SMTP 连接已经被服务端关闭
	if errors.Is(err, io.EOF) {
		return true
	}

	// 网络层错误
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	errStr := strings.ToLower(err.Error())

	// 常见的连接断开错误
	connectionErrors := []string{
		"broken pipe",
		"connection reset by peer",
		"connection reset",
		"closed network connection",
		"use of closed network connection",
		"unexpected eof",
		"connection refused",
		"connection timed out",
		"i/o timeout",
		"no route to host",
	}

	for _, msg := range connectionErrors {
		if strings.Contains(errStr, msg) {
			return true
		}
	}

	return false
}