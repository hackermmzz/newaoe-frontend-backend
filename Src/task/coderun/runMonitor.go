package coderun

import (
	"context"
	api_codeRun "new-aoe-judge/Src/api/codeRun"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/profiler"
	"os"
	"sync"
)

func runMonitor(
	ctx context.Context,
	server grpc_api.CodeClient,
	resource profiler.CPUResourceAllocateInfo,
	id string,
	indices int64,
	runDir string,
	buildDir string,
	crashFile *os.File,
	logFile *os.File,
	recordFile *os.File,
	resultFile *os.File,
	debugFile *os.File,
) (error, *api_codeRun.CodeRunRetInfo) {
	//声明变量
	var runRet *api_codeRun.CodeRunRetInfo
	var streamErr error
	codeRunDone := make(chan struct{}, 1)
	allDone := make(chan struct{}, 1)
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
			return streamErr, nil
		} else if runRet.Err != nil {
			return runRet.Err, nil
		} else {
			return nil, runRet
		}
	case <-ctx.Done():
		return ctx.Err(), nil
	}
}
