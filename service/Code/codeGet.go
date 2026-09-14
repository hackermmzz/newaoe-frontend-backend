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
		util.DebugError("疑似Auth泄露!")
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
	urlGet := dao.OssGetDownloadFileUrls([]string{header_dir, source_dir}, duration, false)
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
	}
	return data, nil
}

func GetOneCodeTask() interface{} {
	session := dao.DB.NewSession()
	defer session.Close()
	var codeinfo dao.CodeRunInfo
	//这里要做幂等，防止这个消息已经被消费过了
	for i := 0; i < 10; i += 1 {
		data, success := dao.RedisListPop(context.Background(), config.Conf.Code.CodeWaitForRunQueueTopic)
		if !success {
			return nil
		}
		err := json.Unmarshal([]byte(data), &codeinfo)
		if err != nil {
			util.DebugError("GetOneCodeTask JsonUnmarshal err:", err)
			return nil
		}
		//判断是否已经处理过了(只有CodeRunning表里面没有，且CodeRun表有才算成功跑结束，取反就是下面这个)
		if dao.CodeRunningExist(session, codeinfo.Indices) || !dao.CodeRunExist(session, codeinfo.Indices) {
			//
			util.DebugSuccess("OJ successfully get one code!")
			return codeinfo
		}
	}
	//
	return nil
}

// 鉴权使用
func codeGrpcServerAuthConfirm(auth string) bool {
	return auth == "5rGq56uL5rSq5piv5YWo5LiW55WM5pyA5biF55qE55S355Sf"
}
