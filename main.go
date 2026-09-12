package main

import (
	"newaoe/config"
	"newaoe/dao"
	"newaoe/service/Grpc"
	"newaoe/service/Server"
	"newaoe/service/Service"
)

func main() {
	//初始化基础配置
	config.ConfigInit()
	//初始化MQ
	dao.RocketMQInit()
	//初始化服务
	Service.ServiceInit()
	//连接数据库
	dao.ConnectDatabase()
	//连接Redis缓存
	dao.ConnectRedis()
	//连接OSS对象存储
	dao.OssInit()
	//启动GRPC
	Grpc.GrpcInit()
	//所有服务初始化完毕后的预处理
	Service.PostProcess()
	//启动后端
	Server.ServeInit()
}
