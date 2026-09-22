package service

import (
	"context"
	"encoding/json"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/redis"
	"newaoe/Src/util"
)

func GetOneCodeTask() (*CodeRunTaskInfo, error) {
	session := database.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return nil, util.NewError(err)
	}
	defer session.Rollback()
	//这里要做幂等，防止这个消息已经被消费过了
	var codeinfo CodeRunTaskInfo
	for i := 0; i < 10; i += 1 {
		data, success := redis.RedisListPop(context.Background(), config.Conf.Code.CodeWaitForRunQueueTopic)
		if !success {
			return nil, nil
		}
		err := json.Unmarshal([]byte(data), &codeinfo)
		if err != nil {
			return nil, util.NewError(err)
		}
		//判断是否已经处理过了(只有CodeRunning表里面没有，且CodeRun表有才算成功跑结束，取反就是下面这个)
		has0, err := dao.CodeRunningExist(session, codeinfo.Indices)
		has1 := dao.CodeRunExist(session, codeinfo.Indices)
		if err != nil {
			return nil, util.NewError(err)
		}
		if has0 || !has1 {
			//
			util.DebugSuccess("OJ successfully get one code!")
			return &codeinfo, nil
		}
	}
	//
	return nil, util.NewError("获取失败!")
}
