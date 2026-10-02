package coderun

import (
	"context"
	"encoding/json"
	"fmt"
	api_codeRun "new-aoe-judge/Src/api/codeRun"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/util"
)

func processCrash(
	ctx context.Context,
	server grpc_api.CodeClient,
	id string,
	indices int64,
	crashFile string,
	debugFile string,
	runRet *api_codeRun.CodeRunRetInfo,
) error {
	//判断是否是崩溃
	needLog := isProgramSelfReason(int32(runRet.ExitCode))
	dataMp := map[string]interface{}{
		"crash_reason": runRet.ExitReason,
		"needlog":      needLog,
	}
	//发送数据给服务器
	resp, err := global.PostCodeStatus(
		ctx,
		server,
		id,
		indices,
		global.Code_Status_Crash,
		global.CodeRunStatusInfo{
			Status: global.Code_Status_Crash,
			Data:   util.JsonMap(dataMp),
		}.String(),
	)
	if err != nil {
		return fmt.Errorf("processCrash:发送崩溃状态失败: %w %v", err, resp.Data)
	}
	//上传崩溃日志和debug日志
	urls := make(map[string]string)
	err = json.Unmarshal([]byte(resp.Data), &urls)
	if err != nil {
		return fmt.Errorf("processCrash:解析崩溃状态更新数据失败: %w %v", err, resp.Data)
	}
	if needLog {
		//上传崩溃日志
		crash_url, exist0 := urls["crash_log_url"]
		if !exist0 {
			return fmt.Errorf("崩溃日志上传链接为空 %v!", resp.Data)
		}
		err = global.UploadFile(crash_url, crashFile)
		if err != nil {
			return fmt.Errorf("processCrash:上传崩溃日志文件失败: %w %v", err, resp.Data)
		}
	}
	//上传debug日志
	stdout_url, exist1 := urls["debug_log_url"]
	if !exist1 {
		return fmt.Errorf("debug日志上传链接为空 %v!", resp.Data)
	}
	err = global.UploadFile(stdout_url, debugFile)
	if err != nil {
		return fmt.Errorf("processCrash:上传debug日志文件失败: %w %v", err, resp.Data)
	}
	return nil
}

func isProgramSelfReason(returnCode int32) bool {
	switch returnCode {
	case 0, 1, 125, 126, 127, 132, 134, 136, 139:
		return true
	case 130, 137, 143:
		return false
	default:
		return true
	}
}
