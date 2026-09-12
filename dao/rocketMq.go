package dao

import (
	"context"
	"newaoe/config"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

var RocketMQProducer rocketmq.Producer

func RocketMQInit() {
	rlog.SetLogLevel("error")
	//初始化Producer(只用一个MQ足以)
	RocketMQProducer = NewMQProducer()
	//
	err := RocketMQProducer.Start()
	if err != nil {
		panic("RocketMQInit:" + err.Error())
	}
	//
}

func NewMQProducer() rocketmq.Producer {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{config.Conf.RocketMQ.Host}),
		producer.WithRetry(3),
	)
	if err != nil {
		panic("NewMQProducer" + err.Error())
	}
	return p
}

func NewMQPullConsumer(group string) rocketmq.PullConsumer {
	c, err := rocketmq.NewPullConsumer(
		consumer.WithNameServer([]string{config.Conf.RocketMQ.Host}),
		consumer.WithGroupName(group),
	)
	if err != nil {
		panic("NewMQPullConsumer" + err.Error())
	}
	return c
}

func NewMQPushConsumer(topic string, process func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error), groups ...string) rocketmq.PushConsumer {
	var groups_ []consumer.Option
	groups_ = append(groups_, consumer.WithNameServer([]string{config.Conf.RocketMQ.Host}))
	for _, group := range groups {
		groups_ = append(groups_, consumer.WithGroupName(group))
	}
	c, err := rocketmq.NewPushConsumer(
		groups_...,
	)
	if err != nil {
		panic("NewMQPushConsumer" + err.Error())
	}
	err = c.Subscribe(topic, consumer.MessageSelector{}, process)
	if err != nil {
		panic("NewMQPushConsumer" + err.Error())
	}
	return c
}

//func MQPull
