package task

import (
	"context"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/task/codecompile"
	"new-aoe-judge/Src/task/codefetch"
	"new-aoe-judge/Src/task/coderun"
)

// Task任务处理任务
func Task(ctx context.Context, server grpc_api.CodeClient) {
	//获取代码任务
	for i := 0; i < global.JudgeMaxCore; i += 1 {
		go codefetch.Task_GetStudentCode(ctx, server)
	}
	//编译任务
	for i := 0; i < global.JudgeMaxCore*100/global.CompileCPULimit; i += 1 {
		go codecompile.Task_ProcessCompile(ctx, server)
	}
	//运行任务
	for i := 0; i < global.JudgeMaxCore*100/global.RunCPULimit; i += 1 {

		go coderun.Task_ProcessRun(ctx, server)
	}
}
