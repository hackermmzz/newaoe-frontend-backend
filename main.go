package main

import (
	"fmt"
	Grpc "newaoe/Src/codeRun/grpc"
	"newaoe/Src/initialize"
	"newaoe/Src/oss"
)

func main() {
	//初始化基础配置
	fmt.Println("正在初始化配置")
	initialize.ConfigInit()
	fmt.Println("配置初始化完成")
	//初始化MQ
	fmt.Println("正在初始化消息队列")
	initialize.MessageQueueInit()
	fmt.Println("消息队列初始化完成")
	//连接数据库
	fmt.Println("正在初始化数据库")
	initialize.DataBaseInit()
	fmt.Println("数据库初始化完成")
	//连接Redis缓存
	fmt.Println("正在初始化redis")
	initialize.RedisInit()
	fmt.Println("redis初始化完成")
	//连接OSS对象存储
	fmt.Println("正在初始化OSS")
	oss.OssInit()
	fmt.Println("OSS初始化完成")
	//初始化邮箱连接
	fmt.Println("正在初始化邮箱服务")
	initialize.EmailInit()
	fmt.Println("邮箱服务初始化完成")
	//初始化代码运行状态服务
	fmt.Println("正在初始化代码运行服务")
	initialize.CodeRunInit()
	fmt.Println("代码运行服务初始化完成")
	//启动GRPC
	fmt.Println("正在初始化GRPC服务")
	Grpc.GrpcInit()
	fmt.Println("GRPC服务初始化完成")
	//所有服务初始化完毕后的预处理
	initialize.PostInit()
	//启动后端
	initialize.ServerInit()
}
