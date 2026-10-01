package codefetch

import (
	"context"
	"fmt"
	api_codeGet "new-aoe-judge/Src/api/codeGet"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/task/codecompile"
	"time"
)

// ////////////////// 获取代码
func Task_GetStudentCode(ctx context.Context, server grpc_api.CodeClient) {
	for {
		func() {
			res := api_codeGet.GetOneStudentCode(ctx, server)
			if !res.OK {
				global.Log(res.Msg)
				time.Sleep(time.Duration(config.Conf.JudgeSleepTimeWhenGetCodeFailed) * time.Second)
				return
			}
			//记录性能
			global.Profiler.IncreaseTaskWait()
			defer global.Profiler.DecreaseTaskWait()
			//
			global.Log(fmt.Sprintf("获取代码成功! ID:%s Indices:%d RunType:%d", res.ID, res.Indices, res.RunType))
			//将代码添加到编译队列
			global.CompileWaitQueue <- codecompile.ForCompileInfo{
				StudentCode: res,
			}
		}()
	}
}
