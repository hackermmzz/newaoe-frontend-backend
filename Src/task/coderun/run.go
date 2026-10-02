package coderun

import (
	"context"
	"fmt"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/util"
	"os"
)

type ForRunInfo struct {
	ID       string
	Indices  int64
	RunDir   string
	BuildDir string
}

// 运行任务
func Task_ProcessRun(ctx context.Context, server grpc_api.CodeClient) {
	for {
		var code ForRunInfo
		var resultFilePath string
		select {
		case <-ctx.Done():
			return
		case dt := <-global.RunWaitQueue:
			code = dt.(ForRunInfo)
		}
		func() {
			//////////////创建文件
			crashLogFilePath := util.JoinPath(code.RunDir, config.Conf.CrashLogFileName)
			logFilePath := util.JoinPath(code.RunDir, config.Conf.RunLogFileName)
			recordFilePath := util.JoinPath(code.RunDir, config.Conf.RecordFileName)
			resultFilePath = util.JoinPath(code.RunDir, config.Conf.RunResultFileName)
			debugLogFilePath := util.JoinPath(code.RunDir, config.Conf.RunDebugLogOutputFileName)
			files, err := util.OpenFiles([]string{
				crashLogFilePath,
				logFilePath,
				recordFilePath,
				resultFilePath,
				debugLogFilePath,
			}, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				global.LogError("创建文件失败", err.Error())
				return
			}
			defer util.CloseFiles(files)
			//记录性能
			global.Profiler.IncreaseRunning()
			defer global.Profiler.DecreaseRunning()
			//打印日志
			global.LogInfo(fmt.Sprintf("%s_%d正在运行...", code.ID, code.Indices))
			///////////////运行代码
			err = runMonitor(
				ctx,
				server,
				code.ID,
				code.Indices,
				code.RunDir,
				code.BuildDir,
				files[0],
				files[1],
				files[2],
				files[3],
				files[4],
			)
			if err != nil {
				global.LogError(fmt.Sprintf("%s/%d/运行代码失败: %s", code.ID, code.Indices, err.Error()))
			} else {
				global.LogSuccess(fmt.Sprintf("%s/%d/运行结束!", code.ID, code.Indices))
			}
		}()
		////////////////////删除日志文件
		if _, err := os.Stat(resultFilePath); err == nil {
			err = os.Remove(resultFilePath)
			if err != nil {
				global.LogError(fmt.Sprintf("%s/%d/删除日志文件失败: %s", code.ID, code.Indices, err.Error()))
			}
		}
	}
}
