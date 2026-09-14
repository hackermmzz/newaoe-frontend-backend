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
func RunUserCode(session *xorm.Session, info dao.CodeRunInfo) bool {
	//向数据库插入一条运行记录
	indices := dao.CodeRunAdd(session, info)
	if indices == 0 {
		util.DebugError("RunUserCode插入记录失败!")
		return false
	}
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
