package coderun

import (
	"context"
	"encoding/json"
	"fmt"
	"new-aoe-judge/Src/config"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/util"
	"os"
	"time"
)

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
