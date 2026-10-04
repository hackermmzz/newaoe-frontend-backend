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
	"os"
	"time"
)

// ////////////////// 获取代码
func Task_GetStudentCode(ctx context.Context, server grpc_api.CodeClient) {
	for {
		//
		select {
		case <-ctx.Done():
			return
		default:
			func() {
				//提前申请一个任务槽
				bNeedComplete := true
				bNeedSleep := true
				global.Profiler.TaskProcess()
				defer func() {
					if bNeedComplete {
						global.Profiler.TaskComplete()
						if bNeedSleep {
							sleepTime := time.Duration(config.Conf.JudgeSleepTimeWhenGetCodeFailed-rand.IntN(2)+1) * time.Second
							time.Sleep(sleepTime)
						}
					}
				}()
				//获取代码
				res := api_codeGet.GetOneStudentCode(ctx, server)

				if res.Error != nil {
					global.LogError(res.Error.Error())
					return
				}
				if !res.OK {
					return
				}
				//处理一下.cpp文件
				err := fixCPP(res.Source)
				if err != nil {
					bNeedSleep = false
					global.LogError(fmt.Sprintf("处理代码失败! %s", err.Error()))
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
					bNeedSleep = false
					return
				case global.CompileWaitQueue <- codecompile.ForCompileInfo{
					StudentCode: res,
				}:
				}
				//保持任务信号量
				bNeedComplete = false
			}()
		}
	}
}

// 目前需要对cpp进行一下处理，增加一个函数就行了
func fixCPP(sourceFile string) error {
	extra_fun := `

		UsrAI* MMZZ_NewUsrAI(){
			return new UsrAI();
		}
			
	`

	// 以追加模式打开文件，O_APPEND 写到文件末尾
	f, err := os.OpenFile(sourceFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(extra_fun)
	return err
}
