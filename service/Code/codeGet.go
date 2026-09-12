package Code

import (
	"context"
	"encoding/json"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/service/Grpc/grpc_api"
	"newaoe/util"
	"time"
)

type GrpcCodeServer struct {
	grpc_api.UnimplementedCodeServer
}

func (server *GrpcCodeServer) GetCode(ctx context.Context, req *grpc_api.CodeRequest) (*grpc_api.CodeReply, error) {
	//鉴权
	if !codeGrpcServerAuthConfirm(req.Auth) {
		util.Debug("疑似Auth泄露!")
		return nil, nil
	}
	//获取一份代码
	info := GetOneCodeTask()
	if info == nil {
		return &grpc_api.CodeReply{
			Ok:  false,
			Msg: "代码队列为空!",
		}, nil
	}
	codeinfo := info.(dao.CodeRunInfo)
	id := codeinfo.ID
	indices := codeinfo.Indices
	source_dir := codeinfo.Source
	header_dir := codeinfo.Header
	//获取链接（这里不走路由了，直接后端获取链接）
	duration := time.Duration(60*30) * time.Second
	headerUrl := dao.OssGetDownloadFileUrl(header_dir, duration, false)
	sourceUrl := dao.OssGetDownloadFileUrl(source_dir, duration, false)
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
	}
	return data, nil
}

func GetOneCodeTask() interface{} {
	data, success := dao.RedisListPop(context.Background(), config.Conf.Code.CodeWaitForRunRedisQueueTopic)
	if !success {
		return nil
	}
	var codeinfo dao.CodeRunInfo
	err := json.Unmarshal([]byte(data), &codeinfo)
	if err != nil {
		util.Debug("GetOneCodeTask JsonUnmarshal err:", err)
		return nil
	}
	util.Debug("OJ successfully get one code!")
	return codeinfo
}

// 鉴权使用
func codeGrpcServerAuthConfirm(auth string) bool {
	return auth == "5rGq56uL5rSq5piv5YWo5LiW55WM5pyA5biF55qE55S355Sf"
}
