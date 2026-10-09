package controller

import (
	"context"
	"newaoe/Src/codeRun/grpc/grpc_api"
	"newaoe/Src/codeRun/model"
	"newaoe/Src/codeRun/service"
	"newaoe/Src/util"
)

func (server *GrpcCodeServer) CodeStatusUpdate(ctx context.Context, req *grpc_api.CodeStatusUpdateRequest) (*grpc_api.StatusUpdateReply, error) {
	//鉴权
	var ok bool
	var auth string
	if ok, auth = service.CodeGrpcServerAuthConfirm(req.Auth); !ok {
		util.DebugError("疑似Auth泄露!")
		return nil, nil
	}
	//获取数据
	indices := req.Indices
	id := req.Id
	info := model.NewCodeRunStatusInfo()
	if info.Unmarshal([]byte(req.Data)) != nil {
		util.DebugError("OJ传过来的数据格式有误!")
		return nil, util.NewError("数据格式错误!请更新oj端!")
	}
	//处理数据
	ret, err := service.CodeRunStatusGetProcess(id, int64(indices), info)
	if err != nil {
		util.DebugError("CodeStatusUpdate", err)
		return nil, util.NewError("数据格式错误!请更新oj端!") //处理失败一定是客户端问题
	}
	if ret == nil {
		ret = &grpc_api.StatusUpdateReply{}
	}
	//
	ret.Auth = auth
	return ret, nil
}
