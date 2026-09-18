package controller

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/grpc/grpc_api"
	"newaoe/Src/codeRun/model"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/mq"
	"newaoe/Src/oss"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"path"
	"strconv"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"xorm.io/xorm"
)

type CodeRunStatusInfoPushRedis struct {
	ID      string `json:"id"`
	Indices int    `json:"indices"`
	model.CodeRunStatusInfo
}

var (
	CodeRunStatusPushConsumer  rocketmq.PushConsumer                                              //代码运行状态消费者
	CodeRunStatusUpdateQueue   = make([]int, 0, config.Conf.Code.CodeRunStatusUpdateQueueMaxSize) //代码运行状态更新队列
	CodeRunStatusUpdateChannel = make(chan int, config.Conf.Code.CodeRunStatusUpdateQueueMaxSize) //代码运行状态更新通道
)

func (server *GrpcCodeServer) CodeStatusUpdate(ctx context.Context, req *grpc_api.CodeStatusUpdateRequest) (*grpc_api.StatusUpdateReply, error) {
	var ret *grpc_api.StatusUpdateReply
	//鉴权
	if !codeGrpcServerAuthConfirm(req.Auth) {
		util.DebugError("疑似Auth泄露!")
		return nil, nil
	}
	//更新数据
	indices := req.Indices
	id := req.Id
	info := model.NewCodeRunStatusInfo()
	if info.Unmarshal([]byte(req.Data)) != nil {
		util.DebugError("OJ传过来的数据格式有误!")
		return nil, nil
	}
	//判断是否需要下载链接(如果是会改变codeRunStatus一些字段)
	ret = processCodeRunStatus(int(indices), id, &info)
	//先更新redis
	redisData := CodeRunStatusInfoPushRedis{
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

// 初始化代码运行监控服务
func CodeRunServiceInit() {
	var groups []string
	for i := 0; i < config.Conf.Code.CodeRunStatusGroupCount; i++ {
		groups = append(groups, "CodeRunStatusGroup"+strconv.Itoa(i))
	}
	//初始化消费者
	CodeRunStatusPushConsumer = mq.NewMQPushConsumer(config.Conf.Code.CodeRunStatusTopic, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		//解析消息

		for _, msg := range msgs {
			indices := int64(binary.BigEndian.Uint64(msg.Body))
			//将获取到的数据放到队列中
			CodeRunStatusUpdateChannel <- int(indices)
		}
		return consumer.ConsumeSuccess, nil
	}, groups...)

	// 启动消费者
	err := CodeRunStatusPushConsumer.Start()
	if err != nil {
		panic("CodeRunServiceInit:" + err.Error())
	}
	//启动队列处理线程
	go codeRunStatusUpdateQueueProcess()
	//启动一个协程定期检查过期没运行的
	go codeRunningLongTimeWaitRepush()
}

func codeRunningLongTimeWaitRepush() {
	for {
		func() {
			session := database.NewSession()
			defer session.Close()
			// 超时时间
			expireDuration := time.Duration(config.Conf.Code.CodeWaitTooLongTimeLimit) * time.Minute
			runningList := dao.CodeRunningGetExpireTime(session, expireDuration, 20)
			if len(runningList) == 0 {
				return
			}
			//获取所有info
			indicesArr := make([]int, len(runningList))
			for i, d := range runningList {
				indicesArr[i] = d.Indices
			}
			coderunInfos := dao.CodeRunBatchGetByIndices(indicesArr)
			for _, info := range coderunInfos {
				if RunUserCode(session, info, false) {
					//成功了就更新当前的重新运行的时间戳
					if !dao.CodeRunningUpdate(session, model.CodeRunningInfo{
						Indices:    info.Indices,
						SubmitTime: util.UTC_Time(),
					}) {
						util.DebugError("CodeRunningUpdate失败!")
					}
				}
			}
			if err := session.Commit(); err != nil {
				util.DebugError("codeRunningLongTimeWaitRepush:session commit fail", err)
			}
		}()
		// 防止CPU空转
		time.Sleep(10 * time.Second)
	}
}
func codeRunStatusUpdateQueueProcess() {
	duration := time.Duration(config.Conf.Code.CodeRunStatusUpdateInterval) * time.Second
	timer := time.NewTicker(duration)
	for {
		select {
		case d := <-CodeRunStatusUpdateChannel:
			CodeRunStatusUpdateQueue = append(CodeRunStatusUpdateQueue, d)
			if len(CodeRunStatusUpdateQueue) >= config.Conf.Code.CodeRunStatusUpdateQueueSizeThreshold {
				batchUpdateCodeRunStatus(CodeRunStatusUpdateQueue)
				CodeRunStatusUpdateQueue = CodeRunStatusUpdateQueue[:0]
				timer.Reset(duration)
			}

		case <-timer.C:
			if len(CodeRunStatusUpdateQueue) > 0 {
				batchUpdateCodeRunStatus(CodeRunStatusUpdateQueue)
				CodeRunStatusUpdateQueue = CodeRunStatusUpdateQueue[:0]
			}
		}
	}
}

// 批量更新代码运行状态
func batchUpdateCodeRunStatus(indices []int) {
	// 1. 对 indices 去重
	distinctIndices := make(map[int]bool)
	for _, ind := range indices {
		distinctIndices[ind] = true
	}
	if len(distinctIndices) == 0 {
		return
	}
	// 2. 从 Redis 获取数据
	ctx := context.Background()
	dataArr := make([]model.CodeRunInfo, len(distinctIndices))
	idx := 0
	version := util.UTC_Time().UnixNano()
	statusIndicesMap := make(map[int]CodeRunStatusInfoPushRedis)
	for ind := range distinctIndices {
		data, exist := redis.RedisGet(ctx, fmt.Sprintf("CodeRun:%v", ind))
		dt := &dataArr[idx]
		dt.Indices = ind
		dt.Version = version
		if exist {
			var v CodeRunStatusInfoPushRedis
			err := json.Unmarshal(data, &v)
			if err != nil {
				util.DebugError("Why err := json.Unmarshal(data, &v) fail?", err)
			}
			dt.Status = v.CodeRunStatusInfo.Marshal()
			statusIndicesMap[ind] = v
		} else {
			status := model.NewCodeRunStatusInfo()
			status.Status = model.Code_Status_Error
			dt.Status = model.ProcessDataMessageByStatus(status).Marshal()
		}
		idx += 1
	}
	// 3. 批量更新到数据库
	// 批量insert临时表
	var err error
	session := database.NewSession()
	defer session.Close()
	//开启事务
	err = session.Begin()
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus:" + err.Error())
		return
	}
	//创建临时表
	temp_table := fmt.Sprintf("tmp_code_run_%v", version)
	_, err = session.Exec(fmt.Sprintf("CREATE TEMPORARY TABLE %v (indices INT, status TEXT, new_version BIGINT)", temp_table))
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: create temp table failed:" + err.Error())
		session.Rollback()
		return
	}
	// 批量insert临时表
	type TempCodeRunInfo struct {
		Indices    int    `xorm:"indices"`
		Status     string `xorm:"text status"`
		NewVersion int64  `xorm:"new_version"`
	}
	tempDataArr := make([]TempCodeRunInfo, len(dataArr))
	for i := range dataArr {
		tempDataArr[i] = TempCodeRunInfo{
			Indices:    dataArr[i].Indices,
			NewVersion: dataArr[i].Version,
		}
		//解析status(它包含完整运行数据)
		tempDataArr[i].Status = processCodeRunningStatus(statusIndicesMap[dataArr[i].Indices])
	}
	_, err = session.Table(temp_table).Insert(tempDataArr)
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: insert temp table failed:" + err.Error())
		session.Rollback()
		return
	}
	// JOIN更新主表
	cmd := fmt.Sprintf(`	
    UPDATE %v c
    JOIN %v t ON c.indices = t.indices
    SET c.status = t.status, c.version = t.new_version
    where c.version <t.new_version
	`, model.CodeRunInfo{}.TableName(), temp_table)
	_, err = session.Exec(cmd)
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: update main table failed:" + err.Error())
		session.Rollback()
		return
	}
	//删除临时表
	_, err = session.Exec(fmt.Sprintf("DROP TEMPORARY TABLE %v", temp_table))
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: drop temp table failed:" + err.Error())
		return
	}
	//更新排行榜
	updateRank(session, dataArr)
	//移除已经结束的记录
	removeCodeRunRecordAlreadyFinish(session, dataArr)
	//提交事务
	err = session.Commit()
	if err != nil {
		util.DebugError("batchUpdateCodeRunStatus: commit transaction failed:", err)
		return
	}
}

func processCodeRunningStatus(status CodeRunStatusInfoPushRedis) string {
	/*这个在推送状态的时候就已经处理了*/
	// runstatus := status.Status
	// data := status.Data
	// //如果是编译失败，那么则存储一个编译失败的下载链接而不是完整的编译失败信息
	// if runstatus == dao.Code_Status_Compile_Fail {
	// 	//生成文件地址
	// 	fileName := fmt.Sprintf("compile_log_%d_%d_.txt", status.Indices, util.UTC_Time().Nanosecond())
	// 	filePath := path.Join(config.Conf.OSS.PrivateBaseFolder, status.ID, config.Conf.User.UserOtherFolder, fileName)
	// 	//把数据推送到服务器
	// 	dao.OssUploadFileData(filePath, []byte(status.Data)) //编译失败得话status.Data就是编译错误信息
	// 	//
	// 	data = filePath
	// }
	// //如果是运行成功/运行失败，那么Data字段存储一个录像下载路径
	// if runstatus == dao.Code_Status_Fail || runstatus == dao.Code_Status_Success {
	// 	//下载链接已经在
	// }
	// //如果是运行崩溃，那么Data字段存一个json
	// if runstatus == dao.Code_Status_Crash {

	// }
	// //
	// status.Data = data
	return status.CodeRunStatusInfo.Marshal()
}

func updateRank(session *xorm.Session, dt []model.CodeRunInfo) {
	//
	type CodeRunMsg struct {
		Score int  `json:"score"`
		Frame int  `json:"frame"`
		Win   bool `json:"win"`
	}
	//
	if len(dt) == 0 {
		return
	}
	//
	RankInfoArr := make([]model.RankInfo, 0)
	var err error
	var info model.CodeRunInfo
	var rankinfo model.RankInfo
	for i := range dt {
		ind := dt[i].Indices
		status := dt[i].Status //这个就是CodeRunStatusInfo
		key := fmt.Sprintf("CodeRunInfoForRank:%v", ind)
		//反序列化数据
		var codeRunStatusInfo model.CodeRunStatusInfo
		err = codeRunStatusInfo.Unmarshal([]byte(status)) //redis的数据肯定对
		if err != nil {
			util.DebugError("Why codeRunStatusInfo.Unmarshal([]byte(status)) fail?")
			continue
		}
		//检测状态码，如果不是在运行就不管
		if codeRunStatusInfo.Status < model.Code_Status_Running {
			continue
		}
		//获取当前状态
		data, exist := redis.RedisGet(context.Background(), key)
		if !exist {
			//从数据库读取
			infoPtr := dao.CodeRunGetByIndices(ind)
			if infoPtr == nil {
				util.DebugError("为什么这里是空指针!")
				continue
			}
			info = *infoPtr
			//写入redis
			data, err = json.Marshal(info)
			redis.RedisSet(context.Background(), key, string(data), time.Duration(30)*time.Minute)
		}
		err = json.Unmarshal(data, &info)
		if err != nil {
			util.DebugError("updateRank: json.Unmarshal failed:" + err.Error())
			continue
		}
		//最终数据
		rankinfo.ID = info.ID
		rankinfo.SubmitTime = info.SubmitTime
		if len(info.Description) == 0 {
			info.Description = "原神启动!"
		}
		msgMp := make(map[string]string)
		msgMp["desc"] = info.Description
		msgMp["status"] = codeRunStatusInfo.Data
		msgByte, _ := json.Marshal(msgMp)
		rankinfo.Msg = string(msgByte)
		rankinfo.Score = codeRunStatusInfo.Score
		rankinfo.Frame = codeRunStatusInfo.Frame
		rankinfo.Win = codeRunStatusInfo.Win
		RankInfoArr = append(RankInfoArr, rankinfo)
	}
	//去重
	distinctID := make(map[string]model.RankInfo)
	finalRankInfo := make([]model.RankInfo, 0)
	for i := range RankInfoArr {
		info := RankInfoArr[i]
		//失败不记录到榜单
		if !info.Win {
			continue
		}
		//
		old, ok := distinctID[info.ID]
		if !ok || dao.RankIsBetter(&info, &old) {
			distinctID[info.ID] = info
		}
	}
	for _, value := range distinctID {
		finalRankInfo = append(finalRankInfo, value)
	}
	//提交更新
	if !dao.RankBatchUpdateOrInsertIfBetter(session, finalRankInfo) {
		util.DebugError("updateRank fail!")
	}
}

func removeCodeRunRecordAlreadyFinish(session *xorm.Session, dt []model.CodeRunInfo) bool {
	needRemove := make([]int, 0)
	for _, info := range dt {
		ind := info.Indices
		status := info.Status
		//解析运行结果
		var codeRunningStatus model.CodeRunStatusInfo
		if err := codeRunningStatus.Unmarshal([]byte(status)); err != nil {
			util.DebugError("status解析失败:", err)
			continue
		}
		//运行状态表明程序没有结束
		if !codeRunningStatus.IsFinish() {
			continue
		}
		//
		needRemove = append(needRemove, ind)
	}
	//批量删除
	if !dao.CodeRunningBatchRemove(session, needRemove) {
		util.DebugError("批量移除CodeRunning失败!执行单个单个删除操作!")
		//单个单个移除
		for _, ind := range needRemove {
			if !dao.CodeRunningRemove(session, ind) {
				util.DebugError("单个删除CodeRunning表记录失败!")
			}
		}
	}
	return true
}
