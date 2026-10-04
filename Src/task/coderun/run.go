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
		select {
		case <-ctx.Done():
			return
		case dt := <-global.RunWaitQueue:
			code = dt.(ForRunInfo)
		}
		func() {
			//减少一个任务信号量(无论如何都需要减少的)
			defer global.Profiler.TaskComplete()
			//
			postError := true
			var postErrorMsg error
			defer func() {
				if postError {
					global.PostServerErrorStatus(ctx, server, code.ID, code.Indices, postErrorMsg)
				}
			}()
			//////////////创建文件
			crashLogFilePath := util.JoinPath(code.RunDir, config.Conf.CrashLogFileName)
			logFilePath := util.JoinPath(code.RunDir, config.Conf.RunLogFileName)
			recordFilePath := util.JoinPath(code.RunDir, config.Conf.RecordFileName)
			resultFilePath := util.JoinPath(code.RunDir, config.Conf.RunResultFileName)
			debugLogFilePath := util.JoinPath(code.RunDir, config.Conf.RunDebugLogOutputFileName)
			files, err := util.OpenFiles([]string{
				crashLogFilePath,
				logFilePath,
				recordFilePath,
				resultFilePath,
				debugLogFilePath,
			}, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				postErrorMsg = fmt.Errorf("创建文件失败: %w", err)
				global.LogError(postErrorMsg.Error())
				return
			}
			defer util.CloseFiles(files)
			//记录性能
			global.Profiler.IncreaseRunning()
			defer global.Profiler.DecreaseRunning()
			//打印日志
			global.LogInfo(fmt.Sprintf("%s_%d正在运行...", code.ID, code.Indices))
			///////////////运行代码
			err, runRet := runMonitor(
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
				if runRet == nil {
					postErrorMsg = fmt.Errorf("运行代码失败: %s", err.Error())
				} else {
					postErrorMsg = fmt.Errorf("运行代码失败: %s,耗时:%d秒", err.Error(), runRet.TimeCost)
				}
				global.LogError(postErrorMsg.Error())
				return
			}
			///////////////刷盘结果
			global.LogInfo(fmt.Sprintf("%s/%d/开始刷盘结果...", code.ID, code.Indices))
			err = util.FlushAllFiles(files)
			if err != nil {
				postErrorMsg = fmt.Errorf("刷盘结果失败: %s,耗时:%d秒", err.Error(), runRet.TimeCost)
				global.LogError(postErrorMsg.Error())
				return
			} else {
				global.LogSuccess(fmt.Sprintf("%s/%d/刷盘结果成功", code.ID, code.Indices))
			}
			///////////////处理结果日志
			err = processCodeRunFinish(
				ctx,
				server,
				code.ID,
				code.Indices,
				resultFilePath,
				crashLogFilePath,
				debugLogFilePath,
				recordFilePath,
				runRet,
			)
			if err != nil {
				postErrorMsg = fmt.Errorf("处理结果日志失败: %s ,耗时:%d秒", err.Error(), runRet.TimeCost)
				global.LogError(postErrorMsg.Error())
				return
			}
			global.LogSuccess(fmt.Sprintf("%s/%d/运行结束,耗时:%d秒", code.ID, code.Indices, runRet.TimeCost))
			//删除AOE二进制文件
			err = util.RemoveFiles([]string{util.JoinPath(code.BuildDir, "newAOE")})
			if err != nil {
				//删除失败,继续运行
				global.LogError(err.Error())
			}
			//取消错误报告
			postError = false
		}()
	}
}
