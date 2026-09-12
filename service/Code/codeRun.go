package Code

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/config"
	"newaoe/dao"
	data "newaoe/service/Data"
	"newaoe/util"
)

// 代码当前的运行状态
var (
	Code_Status_Error           = 0 //服务器异常
	Code_Status_Wait            = 1 //在等待队列里面
	Code_Status_Compile         = 2 //编译中
	Code_Status_Compile_Error   = 3 //编译错误
	Code_Status_Compile_Success = 4 //编译成功
	Code_Status_Running         = 5 //正在运行出结果
	Code_Status_Success         = 6 //运行胜利
	Code_Status_Crash           = 7 //游戏崩溃
	Code_Status_Fail            = 8 //游戏失败

)

// 代码运行状态结构体(存储到数据库)
type CodeRunStatus struct {
	Status int    `json:"status"` //当前状态
	Data   string `json:"data"`   //当前数据
}

func (c CodeRunStatus) String() string {
	res, _ := json.Marshal(c)
	return string(res)
}

// 根据索引将对应的代码文件加入待运行队列
func RunUserCode(indices int) {
	//将代码加入代运行就绪队列
	info := data.CodeRunGetByIndices(indices)
	if info == nil {
		return
	}
	if AddCodeFile(*info) {
		//
		return
	}
}

// 代码文件加入待运行队列
func AddCodeFile(task dao.CodeRunInfo) bool {
	data, err := json.Marshal(task)
	if err != nil {
		util.Debug("AddCodeFile json.Marshal err:", err)
		return false
	}
	return dao.RedisListPush(context.Background(), config.Conf.Code.CodeWaitForRunRedisQueueTopic, string(data))
}

// 判断是否达到提交/运行限制
func LimitCodeUploadOrRun(id string) bool {
	//超级用户不管
	for _, superuser := range config.Conf.Other.SuperUser {
		if superuser == id {
			return false
		}
	}
	//
	key := fmt.Sprintf("CommonUploadOrRunTimes_%v", id)
	time, exist := dao.RedisGet(context.Background(), key)
	curSubmitTime := 0
	if exist {
		err := json.Unmarshal(time, &curSubmitTime)
		if err != nil {
			util.Debug("LimitCodeUploadOrRun Unmarshal Error!", err)
		}
	}
	//判断是否达到限制
	if curSubmitTime >= config.Conf.Code.CodeSubmitTimesPerDay {
		return true
	}
	//
	dao.RedisSet(context.Background(), key, curSubmitTime+1, util.GetLeftTimeForOneDay())
	return false
}
