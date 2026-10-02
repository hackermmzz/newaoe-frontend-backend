package codecompile

import (
	"context"
	"fmt"
	api_codeCompile "new-aoe-judge/Src/api/codeCompile"
	api_codeGet "new-aoe-judge/Src/api/codeGet"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/task/coderun"
	"new-aoe-judge/Src/util"
	"os"
)

type ForCompileInfo struct {
	api_codeGet.StudentCode
}

// 编译任务
func Task_ProcessCompile(ctx context.Context, server grpc_api.CodeClient) {
	for {
		var code ForCompileInfo
		select {
		case <-ctx.Done():
			return
		case dt := <-global.CompileWaitQueue:
			code = dt.(ForCompileInfo)
		}
		//编译代码
		func() {
			//创建编译日志文件
			logFilePath := util.JoinPath(code.BuildDir, config.Conf.CompileLogFileName)
			f, e := os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
			if e != nil {
				global.LogError("创建编译日志文件失败", logFilePath, e.Error())
				return
			}
			defer f.Close()
			//分配资源
			res := global.Profiler.AllocateCodeCompileResource(global.CompileCPULimit)
			defer global.Profiler.RecycleCodeCompileResource(res)
			//记录性能
			global.Profiler.IncreaseCompile()
			defer global.Profiler.DecreaseCompile()
			//通知服务器现在正在编译
			resp, e := global.PostCodeStatus(
				ctx,
				server,
				code.ID,
				code.Indices,
				int32(global.Code_Status_Compile),
				global.CodeRunStatusInfo{
					Status: global.Code_Status_Compile,
				}.String(),
			)
			if e != nil {
				global.LogError("task_ProcessCompile的PostCodeStatus返回错误", e.Error())
				return
			}
			//编译代码
			result := api_codeCompile.CodeCompile(
				res,
				code.ID,
				code.Indices,
				code.BuildDir,
				f,
				code.RunType,
			)
			if result.Error != nil {
				global.LogError("task_ProcessCompile的CodeCompile返回错误", result.Error.Error())
				return
			}
			//处理编译结果
			status := int32(global.Code_Status_Compile_Fail)
			if result.OK {
				status = int32(global.Code_Status_Compile_Success)
			}
			resp, e = global.PostCodeStatus(
				ctx,
				server,
				code.ID,
				code.Indices,
				status,
				global.CodeRunStatusInfo{
					Status: status,
				}.String(),
			)
			if e != nil {
				global.LogError("task_ProcessCompile的PostCodeStatus返回错误", e.Error())
				return
			}
			//处理失败情况
			if !result.OK {
				global.LogSuccess(fmt.Sprintf("%s/%d/编译失败: %s", code.ID, code.Indices, result.Msg))
				if resp == nil {
					global.LogError("task_ProcessCompile的PostCodeStatus返回空指针!")
					return
				}
				urls, err := util.JsonToMap(resp.Data)
				if err != nil {
					global.LogError("解析编译状态更新数据失败", resp.Data, err.Error())
					return
				}
				//上传编译错误日志
				url, ok := urls["compile_error_log_url"].(string)
				if !ok {
					global.LogError("编译状态更新数据中缺少编译错误日志URL", resp.Data)
					return
				}
				if e := global.UploadFile(url, logFilePath); e != nil {
					global.LogError("上传编译错误日志失败", url, e.Error())
					return
				}
				return
			}
			//编译成功
			global.LogSuccess(fmt.Sprintf("%s/%d/编译成功!", code.ID, code.Indices))
			//将编译结果添加到运行队列
			select {
			case <-ctx.Done():
				return
			case global.RunWaitQueue <- coderun.ForRunInfo{
				ID:       code.ID,
				Indices:  code.Indices,
				RunDir:   code.RunDir,
				BuildDir: code.BuildDir,
			}:
			}

		}()
	}
}
