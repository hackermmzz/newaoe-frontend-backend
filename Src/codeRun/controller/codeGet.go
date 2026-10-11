package controller

import (
	"context"
	"newaoe/Src/codeRun/grpc/grpc_api"
	"newaoe/Src/codeRun/service"
	"newaoe/Src/util"
)

type GrpcCodeServer struct {
	grpc_api.UnimplementedCodeServer
}

func (server *GrpcCodeServer) GetCode(ctx context.Context, req *grpc_api.CodeRequest) (*grpc_api.CodeReply, error) {
	//鉴权
	var err error
	var certificate string
	if certificate, err = service.CodeGrpcServerAuthConfirm(req.Auth); err != nil {
		util.DebugError("GetCode:", err)
		return nil, err
	}
	//获取一份代码
	codeinfo, urls, err := service.GetOneCodeTask()
	if err != nil {
		util.DebugError(err)
	}
	if codeinfo == nil || urls == nil {
		return &grpc_api.CodeReply{
			Ok:  false,
			Msg: "代码队列为空!",
		}, nil
	}
	//
	headerUrl := urls[0]
	sourceUrl := urls[1]
	//返回数据
	data := &grpc_api.CodeReply{
		Id:        codeinfo.ID,
		Indices:   int64(codeinfo.Indices),
		SourceUrl: sourceUrl,
		HeaderUrl: headerUrl,
		Ok:        true,
		Msg:       "获取成功!",
		Runtype:   int64(codeinfo.RunType),
		Auth:      certificate,
	}
	return data, nil
}
