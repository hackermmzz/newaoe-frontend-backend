package system

import (
	"context"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"time"
)

func SystemRun(ctx context.Context, server grpc_api.CodeClient) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(config.Conf.SystemInfoFetchInterval) * time.Second):
		}
		func() {
			//获取当前oj最新信息
			info, err := global.GetOJSystemInfo(ctx, server)
			if err != nil {
				global.LogError("获取oj系统信息失败: " + err.Error())
				return
			}
			if info == nil {
				global.LogError("获取oj系统信息失败: info is nil")
				return
			}
			//检查是否有新的版本
			err = systemUpdate(info)
			if err != nil {
				global.LogError("检查新版本失败: " + err.Error())
				return
			}
		}()
	}
}
