package coderun

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	api_codeRun "new-aoe-judge/Src/api/codeRun"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/util"
	"os"
	"strings"
	"sync"
	"time"
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
		case dt := <-global.RunWaitQueue:
			code = dt.(ForRunInfo)
		case <-ctx.Done():
			return
		}
		func() {
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
				global.Log("创建文件失败", err.Error())
				return
			}
			defer util.CloseFiles(files)
			//记录性能
			global.Profiler.IncreaseRunning()
			defer global.Profiler.DecreaseRunning()
			//打印日志
			global.Log("%s_%d正在运行...", code.ID, code.Indices)
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
				global.Log(fmt.Sprintf("%s/%d/运行代码失败: %s", code.ID, code.Indices, err.Error()))
			} else {
				global.Log(fmt.Sprintf("%s/%d/运行结束!", code.ID, code.Indices))
			}
		}()
	}
}

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

// 流式读取结果日志
func streamResults(
	ctx context.Context,
	server grpc_api.CodeClient,
	id string,
	indices int64,
	resultFile string,
	done chan struct{},
) error {
	//重新打开结果日志文件
	f, err := os.OpenFile(resultFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	//缓存
	chunkData := make([]byte, 0)
	readPerbatch := int64(4096)
	//读取结果日志
	for {
		select {
		case <-done:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(config.Conf.CodeRunStatusUploadInterval) * time.Second):
		}
		var latest *global.CodeRunStatusInfo
		//读取文件写入缓存
		extraBytes, err := util.ReadFileLimitFromCurrentOffset(f, readPerbatch)
		if err != nil {
			return err
		}
		//合并缓存
		chunkData = append(chunkData, extraBytes...)
		//解析chunkData,看看能不能提取出最新状态
		latest, chunkData, err = processChunkDataForLatestStatus(chunkData)
		if err != nil {
			return err
		}
		//判断是否存在最新的数据
		if latest == nil {
			continue
		}
		//发送最新状态
		global.PostCodeStatus(
			ctx,
			server,
			id,
			indices,
			global.Code_Status_Running,
			global.CodeRunStatusInfo{
				Status: global.Code_Status_Running,
				Food:   latest.Food,
				Wood:   latest.Wood,
				Gold:   latest.Gold,
				Stone:  latest.Stone,
				Frame:  latest.Frame,
				Win:    latest.Win,
				Score:  latest.Score,
				Data:   latest.String(),
			}.String(),
		)
	}
}

// 解析chunkData,看看能不能提取出最新状态
func processChunkDataForLatestStatus(chunkData []byte) (*global.CodeRunStatusInfo, []byte, error) {
	var ret global.CodeRunStatusInfo
	//解析chunkData
	keyPoint := []int{-1}
	for i, v := range chunkData {
		if v == '\n' {
			//从pre+1-i就为最新的状态
			keyPoint = append(keyPoint, i)
		}
	}
	//取最后2个作为最新状态
	if len(keyPoint) >= 2 {
		pre := keyPoint[len(keyPoint)-2] + 1
		now := keyPoint[len(keyPoint)-1]
		dataBytes := chunkData[pre:now]
		newChunk := chunkData[now+1:]
		//解析dataBytes
		err := json.Unmarshal(dataBytes, &ret)
		if err != nil {
			return nil, chunkData, fmt.Errorf("processChunkDataForLatestStatus:解析最新状态失败(这个可能是bug，按道理不应该出现这种问题的!):%s", err.Error())
		}
		//返回最新状态
		return &ret, newChunk, nil
	}
	//
	return nil, chunkData, nil
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
	//正常结束,打开结果日志文件
	f, err := os.OpenFile(resultFile, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("processCodeRunFinish:打开结果日志文件失败: %w", err)
	}
	defer f.Close()
	//定位到文件末尾
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("processCodeRunFinish:定位到文件末尾失败: %w", err)
	}
	//读取文件内容,直到遇到可以解析的
	var data *global.CodeRunStatusInfo
	needreak := false
	for !needreak {
		line, err := util.ReadLineBack(f)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("processCodeRunFinish:读取结果日志文件失败: %w", err)
			}
			needreak = true
		}
		//解析数据
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var tmp global.CodeRunStatusInfo
		err = json.Unmarshal([]byte(line), &tmp)
		if err == nil {
			data = &tmp
			break
		}
	}
	//如果没有解析出数据,则默认失败
	if data == nil {
		data = &global.CodeRunStatusInfo{
			Status: global.Code_Status_Fail,
			Win:    false,
		}
	}
	//发送最终状态
	win := data.Status == global.Code_Status_Success
	finalStatus := int32(global.Code_Status_Success)
	if !win {
		finalStatus = global.Code_Status_Fail
	}
	resp, err := global.PostCodeStatus(
		ctx,
		server,
		id,
		indices,
		finalStatus,
		global.CodeRunStatusInfo{
			Status: finalStatus,
			Win:    win,
			Score:  data.Score,
			Frame:  data.Frame,
			Food:   data.Food,
			Wood:   data.Wood,
			Gold:   data.Gold,
			Stone:  data.Stone,
			Data:   data.String(),
		}.String(),
	)
	if err != nil {
		return fmt.Errorf("processCodeRunFinish:发送最终状态失败: %w", err)
	}
	//上传录像和日志
	urls := make(map[string]string)
	err = json.Unmarshal([]byte(resp.Data), &urls)
	if err != nil {
		return fmt.Errorf("processCodeRunFinish:解析最终状态更新数据失败: %w %v", err, resp.Data)
	}
	video_url, exist0 := urls["video_url"]
	stdout_url, exist1 := urls["debug_log_url"]
	if !exist0 || !exist1 {
		return fmt.Errorf("processCodeRunFinish:视频url:%s和日志url:%s为空 %v", video_url, stdout_url, resp.Data)
	}
	err = global.UploadFile(video_url, recordFile)
	if err != nil {
		return fmt.Errorf("processCodeRunFinish:上传视频文件失败: %w %v", err, resp.Data)
	}
	err = global.UploadFile(stdout_url, debugFile)
	if err != nil {
		return fmt.Errorf("processCodeRunFinish:上传debug日志文件失败: %w %v", err, resp.Data)
	}
	return nil
}

func processCrash(
	ctx context.Context,
	server grpc_api.CodeClient,
	id string,
	indices int64,
	crashFile string,
	debugFile string,
	runRet *api_codeRun.CodeRunRetInfo,
) error {
	//判断是否是崩溃
	needLog := isProgramSelfReason(int32(runRet.ExitCode))
	dataMp := map[string]interface{}{
		"crash_reason": runRet.ExitReason,
		"needlog":      needLog,
	}
	//发送数据给服务器
	resp, err := global.PostCodeStatus(
		ctx,
		server,
		id,
		indices,
		global.Code_Status_Crash,
		global.CodeRunStatusInfo{
			Status: global.Code_Status_Crash,
			Data:   util.JsonMap(dataMp),
		}.String(),
	)
	if err != nil {
		return fmt.Errorf("processCrash:发送崩溃状态失败: %w %v", err, resp.Data)
	}
	//上传崩溃日志和debug日志
	urls := make(map[string]string)
	err = json.Unmarshal([]byte(resp.Data), &urls)
	if err != nil {
		return fmt.Errorf("processCrash:解析崩溃状态更新数据失败: %w %v", err, resp.Data)
	}
	if needLog {
		//上传崩溃日志
		crash_url, exist0 := urls["crash_log_url"]
		if !exist0 {
			return fmt.Errorf("崩溃日志上传链接为空 %v!", resp.Data)
		}
		err = global.UploadFile(crash_url, crashFile)
		if err != nil {
			return fmt.Errorf("processCrash:上传崩溃日志文件失败: %w %v", err, resp.Data)
		}
	}
	//上传debug日志
	stdout_url, exist1 := urls["debug_log_url"]
	if !exist1 {
		return fmt.Errorf("debug日志上传链接为空 %v!", resp.Data)
	}
	err = global.UploadFile(stdout_url, debugFile)
	if err != nil {
		return fmt.Errorf("processCrash:上传debug日志文件失败: %w %v", err, resp.Data)
	}
	return nil
}

func isProgramSelfReason(returnCode int32) bool {
	switch returnCode {
	case 0, 1, 125, 126, 127, 132, 134, 136, 139:
		return true
	case 130, 137, 143:
		return false
	default:
		return true
	}
}
