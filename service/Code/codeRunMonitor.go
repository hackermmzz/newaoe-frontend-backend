package Code

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/service/Grpc/grpc_api"
	"newaoe/util"
	"strconv"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"xorm.io/xorm"
)

var (
	CodeRunStatusPushConsumer  rocketmq.PushConsumer                                              //代码运行状态消费者
	CodeRunStatusUpdateQueue   = make([]int, 0, config.Conf.Code.CodeRunStatusUpdateQueueMaxSize) //代码运行状态更新队列
	CodeRunStatusUpdateChannel = make(chan int, config.Conf.Code.CodeRunStatusUpdateQueueMaxSize) //代码运行状态更新通道
)

func (server *GrpcCodeServer) CodeStatusUpdate(ctx context.Context, req *grpc_api.CodeStatusUpdateRequest) (*grpc_api.Empty, error) {
	//鉴权
	if !codeGrpcServerAuthConfirm(req.Auth) {
		util.Debug("疑似Auth泄露!")
		return nil, nil
	}
	//更新数据
	indices := req.Indices
	codeRunStatus := ProcessDataMessageByStatus(int(req.Status), req.Data)
	//先更新redis
	ctx1 := context.Background()
	duration := time.Duration(60) * time.Minute                                              //设置60分钟过期
	dao.RedisSet(ctx1, fmt.Sprintf("CodeRun:%v", indices), codeRunStatus.String(), duration) //这里肯定不会乱序，因为judge那边是同步发送的
	//再次push到mq（只push编号）
	indices_byte := make([]byte, 8)
	binary.BigEndian.PutUint64(indices_byte, uint64(indices))
	dao.RocketMQProducer.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			util.Debug("代码状态发送到mq失败:" + err.Error())
		}
	}, primitive.NewMessage(
		config.Conf.Code.CodeRunStatusTopic,
		indices_byte,
	))
	return nil, nil
}

func ProcessDataMessageByStatus(status int, data string) CodeRunStatus {
	var ret CodeRunStatus
	switch status {
	case Code_Status_Wait:
		ret.Data = "排队中..."
		ret.Status = Code_Status_Wait
	case Code_Status_Compile:
		ret.Data = "编译中..."
		ret.Status = Code_Status_Compile
	case Code_Status_Compile_Error:
		ret.Data = "编译错误: " + data
		ret.Status = Code_Status_Compile_Error
	case Code_Status_Compile_Success:
		ret.Data = "编译成功"
		ret.Status = Code_Status_Compile_Success
	case Code_Status_Crash:
		ret.Data = "运行崩溃: " + data
		ret.Status = Code_Status_Crash
	case Code_Status_Fail:
		ret.Data = "游戏失败: " + data
		ret.Status = Code_Status_Fail
	case Code_Status_Success:
		ret.Data = "游戏胜利: " + data
		ret.Status = Code_Status_Success
	case Code_Status_Running:
		ret.Data = "正在奋战: " + data
		ret.Status = Code_Status_Running
	default:
		ret.Data = "服务器异常"
		ret.Status = Code_Status_Error
	}

	return ret
}

// 初始化代码运行监控服务
func CodeRunServiceInit() {
	var groups []string
	for i := 0; i < config.Conf.Code.CodeRunStatusGroupCount; i++ {
		groups = append(groups, "CodeRunStatusGroup"+strconv.Itoa(i))
	}
	//初始化消费者
	CodeRunStatusPushConsumer = dao.NewMQPushConsumer(config.Conf.Code.CodeRunStatusTopic, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
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
	go CodeRunStatusUpdateQueueProcess()
}

func CodeRunStatusUpdateQueueProcess() {
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
	dataArr := make([]dao.CodeRunInfo, len(distinctIndices))
	idx := 0
	version := util.UTC_Time().UnixMilli()
	for ind := range distinctIndices {
		data, exist := dao.RedisGet(ctx, fmt.Sprintf("CodeRun:%v", ind))
		dt := &dataArr[idx]
		dt.Indices = ind
		dt.Version = version
		if exist {
			dt.Status = string(data)
		} else {
			dt.Status = ProcessDataMessageByStatus(Code_Status_Error, "").String()
		}
		idx += 1
	}
	// 3. 批量更新到数据库
	// 批量insert临时表
	var err error
	session := dao.DB.NewSession()
	defer session.Close()
	//开启事务
	err = session.Begin()
	if err != nil {
		util.Debug("batchUpdateCodeRunStatus:" + err.Error())
		return
	}
	//创建临时表
	temp_table := fmt.Sprintf("tmp_code_run_%v", version)
	_, err = session.Exec(fmt.Sprintf("CREATE TEMPORARY TABLE %v (indices INT, status TEXT, new_version BIGINT)", temp_table))
	if err != nil {
		util.Debug("batchUpdateCodeRunStatus: create temp table failed:" + err.Error())
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
			Status:     dataArr[i].Status,
			NewVersion: dataArr[i].Version,
		}
	}
	_, err = session.Table(temp_table).Insert(tempDataArr)
	if err != nil {
		util.Debug("batchUpdateCodeRunStatus: insert temp table failed:" + err.Error())
		session.Rollback()
		return
	}
	// JOIN更新主表
	cmd := fmt.Sprintf(`	
    UPDATE %v c
    JOIN %v t ON c.indices = t.indices
    SET c.status = t.status, c.version = t.new_version
    where c.version <t.new_version
	`, dao.CodeRunInfo{}.TableName(), temp_table)
	_, err = session.Exec(cmd)
	if err != nil {
		util.Debug("batchUpdateCodeRunStatus: update main table failed:" + err.Error())
		session.Rollback()
		return
	}
	//删除临时表
	_, err = session.Exec(fmt.Sprintf("DROP TEMPORARY TABLE %v", temp_table))
	if err != nil {
		util.Debug("batchUpdateCodeRunStatus: drop temp table failed:" + err.Error())
		return
	}
	//更新排行榜
	updateRank(session, dataArr)
	//提交事务
	err = session.Commit()
	if err != nil {
		util.Debug("batchUpdateCodeRunStatus: commit transaction failed:" + err.Error())
		return
	}

}

func updateRank(session *xorm.Session, dt []dao.CodeRunInfo) {
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
	RankInfoArr := make([]dao.RankInfo, 0)
	var err error
	var info dao.CodeRunInfo
	var rankinfo dao.RankInfo
	for i := range dt {
		ind := dt[i].Indices
		status := dt[i].Status
		key := fmt.Sprintf("CodeRunInfoForRank:%v", ind)
		//获取当前状态
		data, exist := dao.RedisGet(context.Background(), key)
		if !exist {
			//从数据库读取
			info = *dao.CodeRunGetByIndices(ind)
			//写入redis
			data, err = json.Marshal(info)
			dao.RedisSet(context.Background(), key, string(data), time.Duration(30)*time.Minute)
		}
		err = json.Unmarshal(data, &info)
		if err != nil {
			util.Debug("updateRank: json.Unmarshal failed:" + err.Error())
			continue
		}
		//反序列化数据
		datajs := make(map[string]interface{})
		err = json.Unmarshal([]byte(status), &datajs)
		if err != nil {
			util.Debug("updateRank: json.Unmarshal datajs failed:" + err.Error())
			continue
		}
		//检测状态码，如果不是在运行就不管
		statusCode, ok := datajs["status"]
		if !ok {
			util.Debug("The status in datajs is an error!", err, status)
			continue
		}
		if statusCode.(int) < Code_Status_Running {
			continue
		}
		//反序列化运行数据
		finaldata, ok := datajs["data"]
		if !ok {
			util.Debug("The data in datajs is an error!", err, status)
			continue
		}
		finaldata = util.GetBracesContent(finaldata.(string))
		var msgInfo CodeRunMsg
		err = json.Unmarshal([]byte(finaldata.(string)), &msgInfo)
		if err != nil || finaldata == "" {
			util.Debug("updateRank: json.Unmarshal failed:"+err.Error(), status)
			continue
		}
		//最终数据
		rankinfo.ID = info.ID
		rankinfo.SubmitTime = info.SubmitTime
		if len(info.Description) > 0 {
			rankinfo.Msg = "描述: " + info.Description + "\n" + "状态: " + status
		} else {
			rankinfo.Msg = status
		}
		rankinfo.Score = msgInfo.Score
		rankinfo.Frame = msgInfo.Frame
		rankinfo.Win = msgInfo.Win
		RankInfoArr = append(RankInfoArr, rankinfo)
	}
	//去重
	distinctID := make(map[string]dao.RankInfo)
	finalRankInfo := make([]dao.RankInfo, 0)
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
		util.Debug("updateRank fail!")
	}
}
