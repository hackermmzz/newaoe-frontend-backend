package main

import (
	Grpc "newaoe/Src/codeRun/grpc"
	"newaoe/Src/initialize"
	"newaoe/Src/oss"
)

func main() {
	//初始化基础配置
	initialize.ConfigInit()
	//初始化MQ
	initialize.MessageQueueInit()
	//连接数据库
	initialize.DataBaseInit()
	//连接Redis缓存
	initialize.RedisInit()
	//连接OSS对象存储
	oss.OssInit()
	//启动GRPC
	Grpc.GrpcInit()
	//初始化邮箱连接
	initialize.EmailInit()
	//初始化代码运行状态服务
	initialize.CodeRunInit()
	//所有服务初始化完毕后的预处理
	initialize.PostInit()
	//启动后端
	initialize.ServerInit()
}
