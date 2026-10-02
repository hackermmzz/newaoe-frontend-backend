package coderun

import (
	"context"
	api_codeRun "new-aoe-judge/Src/api/codeRun"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"os"
	"sync"
)

func runMonitor(
	ctx context.Context,
	server grpc_api.CodeClient,
	id string,
	indices int64,
	runDir string,
	buildDir string,
	crashFile *os.File,
	logFile *os.File,
	recordFile *os.File,
	resultFile *os.File,
	debugFile *os.File,
) error {
	//申请资源
	resource := global.Profiler.AllocateCodeRunResource(global.RunCPULimit)
	defer global.Profiler.RecycleCodeRunResource(resource)
	//声明变量
	var runRet *api_codeRun.CodeRunRetInfo
	var streamErr error
	codeRunDone := make(chan struct{})
	defer close(codeRunDone)
	allDone := make(chan struct{})
	defer close(allDone)
	var wg sync.WaitGroup
	wg.Add(2)
	//异步运行代码
	go func() {
		defer wg.Done()
		runRet = api_codeRun.CodeRun(
			resource,
			id,
			indices,
			runDir,
			buildDir,
			crashFile,
			logFile,
			recordFile,
			resultFile,
			debugFile,
		)
		codeRunDone <- struct{}{}
	}()
	//异步读取结果日志
	go func() {
		defer wg.Done()
		streamErr = streamResults(ctx, server, id, indices, resultFile.Name(), codeRunDone)
	}()
	//等待2个人都结束
	go func() {
		wg.Wait()
		allDone <- struct{}{}
	}()
	//等待结果
	select {
	case <-allDone:
		if streamErr != nil {
			return streamErr
		} else if runRet.Err != nil {
			return runRet.Err
		} else {
			return processCodeRunFinish(
				ctx,
				server,
				id,
				indices,
				resultFile.Name(),
				crashFile.Name(),
				debugFile.Name(),
				recordFile.Name(),
				runRet,
			)
		}
	case <-ctx.Done():
		return ctx.Err()
	}
}

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
	//判断是否是崩溃
	if runRet.ExitCode != 0 {
		return processCrash(ctx, server, id, indices, crashFile, debugFile, runRet)
	}
	//处理正常结束
	return processNormalEnd(ctx, server, id, indices, resultFile, recordFile, debugFile)
}
