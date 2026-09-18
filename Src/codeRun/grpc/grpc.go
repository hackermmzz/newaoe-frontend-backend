package Grpc

import (
	"net"
	"newaoe/Src/codeRun/controller"
	"newaoe/Src/codeRun/grpc/grpc_api"
	"newaoe/Src/config"
	"newaoe/Src/util"
	"strconv"

	"google.golang.org/grpc"
)

func GrpcInit() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.Conf.GRPC.Port))
	if err != nil {
		panic("Grpc服务初始化失败!" + err.Error())
	}
	//监听端口|注册服务
	s := grpc.NewServer()
	RegistServer(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			panic("Grpc服务初始化失败!" + err.Error())
		}
	}()
	//
	util.DebugSuccess("GRPC服务启动成功!")
}

func RegistServer(server *grpc.Server) {
	grpc_api.RegisterCodeServer(server, &controller.GrpcCodeServer{})
}
