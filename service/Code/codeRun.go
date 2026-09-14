package Code

import (
	"context"
	"encoding/json"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"

	"xorm.io/xorm"
)

// 根据索引将对应的代码文件加入待运行队列
func RunUserCode(session *xorm.Session, info dao.CodeRunInfo, needInsert bool) bool {
	if needInsert {
		//向数据库插入一条运行记录
		indices := dao.CodeRunAdd(session, info)
		if indices == 0 {
			util.DebugError("RunUserCode插入CodeRunAdd记录失败!")
			return false
		}
		//插入一条记录（以防止崩溃可以恢复）
		if !dao.CodeRunningInsert(session, dao.CodeRunningInfo{Indices: indices, SubmitTime: info.SubmitTime}) {
			util.DebugError("RunUserCode插入CodeRunningInsert记录失败!")
			return false
		}
	} else {
		//清空原先状态
		status := dao.NewCodeRunStatusInfo()
		if !dao.CodeRunUpdateStatusAsync(session, info.Indices, status.Marshal()) {
			util.DebugError("RunUserCode update status fail!")
			return false
		}
	}
	//插入队列，准备给oj消费
	if !addCodeFile(info) {
		//
		return false
	}
	return true
}

// 代码文件加入待运行队列
func addCodeFile(task dao.CodeRunInfo) bool {
	data, err := json.Marshal(task)
	if err != nil {
		util.Debug("AddCodeFile json.Marshal err:", err)
		return false
	}
	return dao.RedisListPush(context.Background(), config.Conf.Code.CodeWaitForRunQueueTopic, string(data))
}
