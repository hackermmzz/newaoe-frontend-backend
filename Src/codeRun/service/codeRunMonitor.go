package service

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"newaoe/Src/codeRun/grpc/grpc_api"
	"newaoe/Src/codeRun/model"
	"newaoe/Src/config"
	"newaoe/Src/mq"
	"newaoe/Src/oss"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"path"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type codeRunStatusInfoPushRedis struct {
	ID      string `json:"id"`
	Indices int    `json:"indices"`
	model.CodeRunStatusInfo
}

func CodeRunStatusGetProcess(id string, indices int64, info model.CodeRunStatusInfo) (*grpc_api.StatusUpdateReply, error) {
	//判断是否需要下载链接(如果是会改变codeRunStatus一些字段)
	ret := processCodeRunStatus(int(indices), id, &info)
	//先更新redis
	redisData := codeRunStatusInfoPushRedis{
		ID:                id,
		Indices:           int(indices),
		CodeRunStatusInfo: info,
	}
	byte_data, _ := json.Marshal(redisData)
	ctx1 := context.Background()
	duration := time.Duration(60) * time.Minute                                   //设置60分钟过期
	redis.RedisSet(ctx1, fmt.Sprintf("CodeRun:%v", indices), byte_data, duration) //这里肯定不会乱序，因为judge那边是同步发送的
	//再次push到mq（只push编号）
	indices_byte := make([]byte, 8)
	binary.BigEndian.PutUint64(indices_byte, uint64(indices))
	msg := primitive.NewMessage(config.Conf.Code.CodeRunStatusTopic, indices_byte)
	pushCodeRunStatusNeedUpdateToQueue(msg, 3) //默认重试3次
	//
	return ret, nil
}

func processCodeRunStatus(indices int, id string, codeRunstatus *model.CodeRunStatusInfo) *grpc_api.StatusUpdateReply {
	status := codeRunstatus.Status
	fp := ""
	dtMap := make(map[string]interface{})
	switch status {
	case model.Code_Status_Compile_Fail:
		fileName := fmt.Sprintf("compile_%d_%d.log", indices, util.UTC_Time().Nanosecond())
		fp = path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserOtherFolder, fileName)
		dtMap["compile_error_log"] = fp
	case model.Code_Status_Crash:
		//解析原来的data，needlog为一个bool，表示是否生成崩溃日志
		mp := make(map[string]interface{})
		json.Unmarshal([]byte(codeRunstatus.Data), &mp)
		if needlog, ok := mp["needlog"].(bool); ok && needlog {
			fileName := fmt.Sprintf("crash_%d_%d.log", indices, util.UTC_Time().Nanosecond())
			fp = path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCrashFolder, fileName)
		}
		if crash_reason, ok := mp["crash_reason"].(string); ok {
			dtMap["crash_reason"] = crash_reason
		} else {
			dtMap["crash_reason"] = ""
		}
		dtMap["crash_log_file"] = fp
	case model.Code_Status_Fail:
		fileName := fmt.Sprintf("video_fail_%d_%d.video", indices, util.UTC_Time().Nanosecond())
		fp = path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserVideoFolder, fileName)
		dtMap["video_file"] = fp
	case model.Code_Status_Success:
		fileName := fmt.Sprintf("video_success_%d_%d.video", indices, util.UTC_Time().Nanosecond())
		fp = path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserVideoFolder, fileName)
		dtMap["video_file"] = fp
	}
	//无论链接生成成功与否，都不用管
	if len(dtMap) != 0 {
		bytes, _ := json.Marshal(dtMap)
		codeRunstatus.Data = string(bytes)
	}
	if fp != "" {
		//尝试3次
		for i := 0; i < 3; i += 1 {
			expire_dur := time.Duration(60) * time.Minute
			url := oss.GetUploadFileUrls([]string{fp}, []time.Duration{expire_dur})
			if len(url) == 1 {
				return &grpc_api.StatusUpdateReply{Data: url[0]}
			}
		}
	}
	//
	return nil
}

func pushCodeRunStatusNeedUpdateToQueue(msg *primitive.Message, limit int) {
	if limit <= 0 {
		util.DebugError("代码运行状态推送队列超出最大重试次数!")
		return
	}
	mq.RocketMQProducer.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			util.DebugError("代码状态发送到mq失败:" + err.Error())
			pushCodeRunStatusNeedUpdateToQueue(msg, limit)
		}
	}, msg)
}
