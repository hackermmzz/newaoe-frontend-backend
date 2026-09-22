package controller

import (
	"context"
	"newaoe/Src/codeRun/grpc/grpc_api"
	"newaoe/Src/codeRun/service"
	"newaoe/Src/oss"
	"newaoe/Src/util"
	"time"
)

type GrpcCodeServer struct {
	grpc_api.UnimplementedCodeServer
}

func (server *GrpcCodeServer) GetCode(ctx context.Context, req *grpc_api.CodeRequest) (*grpc_api.CodeReply, error) {
	//鉴权
	if !codeGrpcServerAuthConfirm(req.Auth) {
		util.DebugError("疑似Auth泄露!")
		return nil, nil
	}
	//获取一份代码
	codeinfo, err := service.GetOneCodeTask()
	if err != nil {
		util.DebugError(err)
	}
	if codeinfo == nil {
		return &grpc_api.CodeReply{
			Ok:  false,
			Msg: "代码队列为空!",
		}, nil
	}
	id := codeinfo.ID
	indices := codeinfo.Indices
	source_dir := codeinfo.Source
	header_dir := codeinfo.Header
	//获取链接（这里不走路由了，直接后端获取链接）
	duration := time.Duration(60*30) * time.Second
	urlGet := oss.OssGetDownloadFileUrls([]string{header_dir, source_dir}, duration, false)
	headerUrl := urlGet[0]
	sourceUrl := urlGet[1]
	if headerUrl == "" || sourceUrl == "" {
		return &grpc_api.CodeReply{
			Ok:  false,
			Msg: "服务器异常!",
		}, nil
	}
	//返回数据
	data := &grpc_api.CodeReply{
		Id:        id,
		Indices:   int64(indices),
		SourceUrl: sourceUrl,
		HeaderUrl: headerUrl,
		Ok:        true,
		Msg:       "获取成功!",
		Runtype:   int64(codeinfo.RunType),
	}
	return data, nil
}

// 鉴权使用
func codeGrpcServerAuthConfirm(auth string) bool {
	return auth == "5rGq56uL5rSq5piv5YWo5LiW55WM5pyA5biF55qE55S355Sf"
}
