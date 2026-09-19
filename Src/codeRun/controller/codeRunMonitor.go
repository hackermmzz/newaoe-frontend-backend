package controller

import (
	"context"
	"encoding/binary"
	"newaoe/Src/codeRun/grpc/grpc_api"
	"newaoe/Src/codeRun/model"
	"newaoe/Src/codeRun/service"
	"newaoe/Src/config"
	"newaoe/Src/mq"
	"newaoe/Src/util"
	"strconv"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

var (
	CodeRunStatusPushConsumer  rocketmq.PushConsumer                                              //代码运行状态消费者
	CodeRunStatusUpdateQueue   = make([]int, 0, config.Conf.Code.CodeRunStatusUpdateQueueMaxSize) //代码运行状态更新队列
	CodeRunStatusUpdateChannel = make(chan int, config.Conf.Code.CodeRunStatusUpdateQueueMaxSize) //代码运行状态更新通道
)

func (server *GrpcCodeServer) CodeStatusUpdate(ctx context.Context, req *grpc_api.CodeStatusUpdateRequest) (*grpc_api.StatusUpdateReply, error) {
	//鉴权
	if !codeGrpcServerAuthConfirm(req.Auth) {
		util.DebugError("疑似Auth泄露!")
		return nil, nil
	}
	//获取数据
	indices := req.Indices
	id := req.Id
	info := model.NewCodeRunStatusInfo()
	if info.Unmarshal([]byte(req.Data)) != nil {
		util.DebugError("OJ传过来的数据格式有误!")
		return nil, nil
	}
	//处理数据
	ret, err := service.CodeRunStatusGetProcess(id, int64(indices), info)
	if err != nil {
		util.DebugError("CodeStatusUpdate", err)
	}
	//
	return ret, nil
}

// 初始化代码运行监控服务
func CodeRunServiceInit() {
	var groups []string
	for i := 0; i < config.Conf.Code.CodeRunStatusGroupCount; i++ {
		groups = append(groups, "CodeRunStatusGroup"+strconv.Itoa(i))
	}
	//初始化消费者
	CodeRunStatusPushConsumer = mq.NewMQPushConsumer(config.Conf.Code.CodeRunStatusTopic, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		//解析消息

		for _, msg := range msgs {
			indices := int64(binary.BigEndian.Uint64(msg.Body))
			//将获取到的数据放到队列中
			CodeRunStatusUpdateChannel <- int(indices)
		}
		return consumer.ConsumeSuccess, nil
	}, groups...)

	// 启动消费者
	err := CodeRunStatusPushConsumer.Start()
	if err != nil {
		panic("CodeRunServiceInit:" + err.Error())
	}
	//启动队列处理线程
	go codeRunStatusUpdateQueueProcess()
	//启动一个协程定期检查过期没运行的
	go codeRunningLongTimeWaitRepush()
}

func codeRunningLongTimeWaitRepush() {
	for {
		func() {
			err := service.ReRunLongTimeWaitRecord()
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
	for {
		select {
		case d := <-CodeRunStatusUpdateChannel:
			CodeRunStatusUpdateQueue = append(CodeRunStatusUpdateQueue, d)
			if len(CodeRunStatusUpdateQueue) >= config.Conf.Code.CodeRunStatusUpdateQueueSizeThreshold {
				service.BatchUpdateCodeRunStatus(CodeRunStatusUpdateQueue)
				CodeRunStatusUpdateQueue = CodeRunStatusUpdateQueue[:0]
				timer.Reset(duration)
			}

		case <-timer.C:
			if len(CodeRunStatusUpdateQueue) > 0 {
				service.BatchUpdateCodeRunStatus(CodeRunStatusUpdateQueue)
				CodeRunStatusUpdateQueue = CodeRunStatusUpdateQueue[:0]
			}
		}
	}
}
