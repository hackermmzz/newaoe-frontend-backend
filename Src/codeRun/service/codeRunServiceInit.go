package service

import (
	"encoding/binary"
	"errors"
	"newaoe/Src/config"
	"newaoe/Src/mq"
	"newaoe/Src/util"
	"time"

	"github.com/apache/rocketmq-clients/golang/v5"
)

var (
	//等待运行
	CodeWaitForRunPullConsumer golang.SimpleConsumer
	//状态更新
	CodeRunStatusPushConsumer  golang.PushConsumer //代码运行状态消费者
	CodeRunStatusUpdateChannel chan int            //代码运行状态更新通道
)

// 初始化代码运行监控服务
func CodeRunServiceInit() {
	//初始化通道(队列大小为阈值的2倍)
	CodeRunStatusUpdateChannel = make(chan int, config.Conf.Code.CodeRunStatusUpdateQueueSizeThreshold*2)
	//初始化状态更新消费者
	err := codeRunStatusConsumerInit()
	if err != nil {
		panic(err.Error())
	}
	//初始化等待运行消费者
	err = codeWaitForRunConsumerInit()
	if err != nil {
		panic(err.Error())
	}
	//启动队列处理线程
	for i := 0; i < config.Conf.Code.CodeRunStatusUpdateCoroutineCount; i += 1 {
		go codeRunStatusUpdateQueueProcess()
	}
	//启动一个协程定期检查过期没运行的
	go codeRunningLongTimeWaitRepush()
}

func codeRunStatusConsumerInit() error {
	//初始化消费者
	CodeRunStatusPushConsumer = mq.NewMQPushConsumer(
		config.Conf.Code.CodeRunStatusTopic,
		func(msg *golang.MessageView) golang.ConsumerResult {
			// 解析消息
			body := msg.GetBody()

			if len(body) < 8 {
				util.DebugError("CodeRunStatus消息长度不足8字节")
				return golang.FAILURE
			}

			indices := int64(binary.BigEndian.Uint64(body))
			// 将获取到的数据放到队列中
			CodeRunStatusUpdateChannel <- int(indices)
			return golang.SUCCESS
		},
		"CodeRunStatus_group",
	)
	// 启动消费者
	err := CodeRunStatusPushConsumer.Start()
	if err != nil {
		return errors.New("CodeRunServiceInit:" + err.Error())
	}
	return nil
}

func codeWaitForRunConsumerInit() error {
	CodeWaitForRunPullConsumer = mq.NewMQSimpleConsumer(
		config.Conf.Code.CodeWaitForRunQueueTopic,
		"CodeWaitForRunGroup",
	)
	if err := CodeWaitForRunPullConsumer.Start(); err != nil {
		return errors.New("codeWaitForRunConsumerInit:" + err.Error())
	}
	return nil
}

func codeRunningLongTimeWaitRepush() {
	for {
		func() {
			expireDuration := time.Duration(config.Conf.Code.CodeWaitTooLongTimeLimit) * time.Minute
			_, err := ReRunLongTimeWaitRecord(expireDuration)
			if err != nil {
				util.DebugError("codeRunningLongTimeWaitRepush:", err)
			}
		}()
		// 防止CPU空转
		time.Sleep(10 * time.Second)
	}
}

func codeRunStatusUpdateQueueProcess() {
	duration := time.Duration(config.Conf.Code.CodeRunStatusUpdateInterval) * time.Second
	timer := time.NewTicker(duration)
	//初始化Queue
	CodeRunStatusUpdateQueue := make([]int, 0, config.Conf.Code.CodeRunStatusUpdateQueueSizeThreshold)
	// 批量更新函数
	fn := func() error {
		infos, err := BatchUpdateCodeRunStatus(CodeRunStatusUpdateQueue)
		if err != nil {
			return util.NewError("BatchUpdateCodeRunStatus:", err)
		}
		err = UpdateRank(infos)
		if err != nil {
			return util.NewError("UpdateRank:", err)
		}
		return nil
	}
	// 处理队列
	for {
		var err error
		select {
		case d := <-CodeRunStatusUpdateChannel:
			CodeRunStatusUpdateQueue = append(CodeRunStatusUpdateQueue, d)
			if len(CodeRunStatusUpdateQueue) >= config.Conf.Code.CodeRunStatusUpdateQueueSizeThreshold {
				err = fn()
				CodeRunStatusUpdateQueue = CodeRunStatusUpdateQueue[:0]
				timer.Reset(duration)
			}

		case <-timer.C:
			if len(CodeRunStatusUpdateQueue) > 0 {
				err = fn()
				CodeRunStatusUpdateQueue = CodeRunStatusUpdateQueue[:0]
			}
		}
		// 处理错误
		if err != nil {
			util.DebugError("codeRunStatusUpdateQueueProcess:", err)
			continue
		}
	}
}
