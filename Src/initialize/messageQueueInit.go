package initialize

import "newaoe/Src/mq"

func MessageQueueInit() {
	//初始化MQ
	mq.RocketMQInit()
}
