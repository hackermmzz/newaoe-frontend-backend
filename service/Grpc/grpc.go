package Grpc

import (
	"net"
	"newaoe/config"
	"newaoe/service/Code"
	"newaoe/service/Grpc/grpc_api"
	"newaoe/util"
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
	util.Debug("GRPC服务启动成功!")
}

func RegistServer(server *grpc.Server) {
	grpc_api.RegisterCodeServer(server, &Code.GrpcCodeServer{})
}
