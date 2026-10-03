package coderun

import (
	"context"
	"fmt"
	api_codeRun "new-aoe-judge/Src/api/codeRun"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
)

// 处理运行结束
func processCodeRunFinish(
	ctx context.Context,
	server grpc_api.CodeClient,
	id string,
	indices int64,
	resultFile string,
	crashFile string,
	debugFile string,
	recordFile string,
	runRet *api_codeRun.CodeRunRetInfo,
) error {
	//如果游戏胜利直接走正常结束路径
	var resultData *global.CodeRunStatusInfo
	var success bool
	var err error
	resultData, success, err = checkIfWin(resultFile)
	if err != nil {
		return fmt.Errorf("检查结果失败: %w", err)
	}
	if success {
		return processNormalEnd(ctx, server, id, indices, resultData, resultFile, recordFile, debugFile, runRet)
	}
	//判断是否是崩溃
	if runRet.ExitCode != 0 {
		return processCrash(ctx, server, id, indices, crashFile, debugFile, runRet)
	}
	//处理正常结束
	return processNormalEnd(ctx, server, id, indices, resultData, resultFile, recordFile, debugFile, runRet)
}

// 获取结果日志最后一行，如果是胜利那么无论崩溃与否直接按照胜利来算
func checkIfWin(resultFile string) (*global.CodeRunStatusInfo, bool, error) {
	//获取最终结果
	res, err := getFinalResult(resultFile)
	if err != nil {
		return nil, false, err
	}
	return res, res.Win, nil
}
