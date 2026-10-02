package coderun

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"new-aoe-judge/Src/global"
	grpc_api "new-aoe-judge/Src/grpc"
	"new-aoe-judge/Src/util"
	"os"
	"strings"
)

func processNormalEnd(
	ctx context.Context,
	server grpc_api.CodeClient,
	id string,
	indices int64,
	resultFile string,
	recordFile string,
	debugFile string,
) error {
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
