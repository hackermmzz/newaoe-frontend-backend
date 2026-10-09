package mq

import (
	"context"
	"newaoe/Src/config"
	"os"
	"time"

	rocketmq "github.com/apache/rocketmq-clients/golang/v5"
)

var RocketMQProducer rocketmq.Producer

func RocketMQInit() {
	// 关闭控制台日志
	_ = os.Setenv("mq.consoleAppender.enabled", "false")
	rocketmq.ResetLogger()

	// 如果你的 Proxy 没启用 TLS
	rocketmq.EnableSsl = false

	// 初始化 Producer
	RocketMQProducer = NewMQProducer()

	err := RocketMQProducer.Start()
	if err != nil {
		panic("RocketMQInit:" + err.Error())
	}
}

// 创建 Producer
func NewMQProducer() rocketmq.Producer {
	p, err := rocketmq.NewProducer(
		&rocketmq.Config{
			// 注意：这里不再是 NameServer :9876
			// 要填 RocketMQ Proxy gRPC 地址，比如 127.0.0.1:8081
			Endpoint: config.Conf.RocketMQ.Host,
		},
		rocketmq.WithTopics(
			config.Conf.Email.EmailMQTopic, //这里不加topic就会一直卡在start阶段，这是框架的bug
			// 其他 topic...
		),
	)
	if err != nil {
		panic("NewMQProducer:" + err.Error())
	}

	return p
}

// 创建 PushConsumer
func NewMQPushConsumer(topic string, process func(*rocketmq.MessageView) rocketmq.ConsumerResult, group string) rocketmq.PushConsumer {
	c, err := rocketmq.NewPushConsumer(
		&rocketmq.Config{
			Endpoint:      config.Conf.RocketMQ.Host,
			ConsumerGroup: group,
		},

		// 订阅 Topic
		rocketmq.WithPushSubscriptionExpressions(
			map[string]*rocketmq.FilterExpression{
				topic: rocketmq.SUB_ALL,
			},
		),

		// 消费函数
		rocketmq.WithPushMessageListener(
			&rocketmq.FuncMessageListener{
				Consume: process,
			},
		),

		// 长轮询等待时间
		rocketmq.WithPushAwaitDuration(
			5*time.Second,
		),
	)

	if err != nil {
		panic("NewMQPushConsumer:" + err.Error())
	}

	return c
}

// 创建 SimpleConsumer
//
// SimpleConsumer 就是你需要的“主动拉取”模式。
// 不会自己不断把任务塞给你，只有你主动调用 Receive() 才拿消息。
func NewMQSimpleConsumer(topic string, group string) rocketmq.SimpleConsumer {
	c, err := rocketmq.NewSimpleConsumer(
		&rocketmq.Config{
			Endpoint:      config.Conf.RocketMQ.Host,
			ConsumerGroup: group,
		},

		// Receive 没消息的时候最长等待多久
		rocketmq.WithSimpleAwaitDuration(
			5*time.Second,
		),

		// 订阅 Topic
		rocketmq.WithSimpleSubscriptionExpressions(
			map[string]*rocketmq.FilterExpression{
				topic: rocketmq.SUB_ALL,
			},
		),
	)

	if err != nil {
		panic("NewMQSimpleConsumer:" + err.Error())
	}

	return c
}

// ACK 消息
func MQAck(c rocketmq.SimpleConsumer, msg *rocketmq.MessageView) error {
	return c.Ack(
		context.Background(),
		msg,
	)
}
