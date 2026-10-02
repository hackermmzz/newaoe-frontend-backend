package codefetch

import (
	"context"
	"fmt"
	"math/rand/v2"
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
		select {
		case <-ctx.Done():
			return
		default:
			func() {
				res := api_codeGet.GetOneStudentCode(ctx, server)
				sleepTime := time.Duration(config.Conf.JudgeSleepTimeWhenGetCodeFailed-rand.IntN(2)+1) * time.Second
				if res.Error != nil {
					global.LogError(res.Error.Error())
					time.Sleep(sleepTime)
					return
				}
				if !res.OK {
					time.Sleep(sleepTime)
					return
				}
				//记录性能
				global.Profiler.IncreaseTaskWait()
				defer global.Profiler.DecreaseTaskWait()
				//
				global.LogSuccess(fmt.Sprintf("获取代码成功! ID:%s Indices:%d RunType:%d", res.ID, res.Indices, res.RunType))
				//将代码添加到编译队列
				select {
				case <-ctx.Done():
					return
				case global.CompileWaitQueue <- codecompile.ForCompileInfo{
					StudentCode: res,
				}:
				}
			}()
		}
	}
}
