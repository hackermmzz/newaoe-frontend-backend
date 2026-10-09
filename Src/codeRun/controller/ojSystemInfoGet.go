package controller

import (
	"context"
	"newaoe/Src/codeRun/grpc/grpc_api"
	"newaoe/Src/codeRun/service"
)

func (server *GrpcCodeServer) OJSystemInfoGet(ctx context.Context, in *grpc_api.Empty) (*grpc_api.OJSystemInfoReply, error) {
	return service.GetOJSystemInfo()
}
