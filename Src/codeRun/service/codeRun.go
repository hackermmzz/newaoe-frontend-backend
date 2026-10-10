package service

import (
	"context"
	"encoding/json"
	"errors"
	codeAssessmentDao "newaoe/Src/codeAssessmentSubmit/dao"
	codeCommonDao "newaoe/Src/codeCommonSubmit/dao"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/codeRun/model"
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/mq"
	"newaoe/Src/util"

	rocketmq "github.com/apache/rocketmq-clients/golang/v5"
)

type CodeRunTaskInfo struct {
	model.CodeRunInfo
	RunType int `json:"runtype"`
}

/*
第一次运行代码时候执行
indices表示提交的索引(服务会根据class判断去哪个表查找)
*/
func CodeRun(indices int, id string, class int, runType int) error {
	//获取代码运行必要的结构数据
	codeRunInfo, err := getCodeRunMustInfo(id, indices, class)
	if err != nil {
		return err
	}
	//
	session := database.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return errors.New("开启事务失败!")
	}
	defer session.Rollback()
	//向数据库插入一条运行记录
	codeRunInfo.Indices, err = dao.CodeRunAdd(session, codeRunInfo)
	if err != nil {
		return err
	}
	if codeRunInfo.Indices == 0 {
		return errors.New("插入CodeRunAdd记录失败!")
	}
	//插入一条记录（以防止崩溃可以恢复）
	inserted, err := dao.CodeRunningInsert(session, model.CodeRunningInfo{
		Indices:    codeRunInfo.Indices,
		SubmitTime: codeRunInfo.SubmitTime,
		RunType:    runType,
	})
	if err != nil {
		return err
	}
	if !inserted {
		return errors.New("插入CodeRunningInsert记录失败!")
	}
	//提交事务
	if err := session.Commit(); err != nil {
		return errors.New("commit失败!")
	}
	//插入队列，准备给oj消费
	if err = addCodeFile(CodeRunTaskInfo{
		CodeRunInfo: codeRunInfo,
		RunType:     runType,
	}, 3); err != nil {
		return err
	}
	//
	return nil
}

// 根据代码提交类别返回运行需要的数据
func getCodeRunMustInfo(id string, indices int, class int) (model.CodeRunInfo, error) {
	codeRunInfo := model.CodeRunInfo{
		ID:         id,
		SubmitTime: util.UTC_Time(),
		Class:      class,
		Status:     model.ProcessDataMessageByStatus(model.NewCodeRunStatusInfo()).Marshal(),
	}
	//分类处理
	switch class {
	case model.Code_CommonSubmit:
		//查表拿到数据
		data, err := codeCommonDao.CodeCommonGetByIndices(nil, indices)
		if err != nil {
			return codeRunInfo, err
		}
		if data == nil || data.ID != id {
			return codeRunInfo, errors.New("有人想作假:" + id)
		}
		codeRunInfo.Header = data.Header
		codeRunInfo.Source = data.Source
		codeRunInfo.Description = data.Description
	case model.Code_AssessmentSubmit:
		data, err := codeAssessmentDao.CodeAssessmentGetByIndices(nil, indices)
		if err != nil {
			return codeRunInfo, err
		}
		if data == nil || data.ID != id {
			return codeRunInfo, errors.New("有人想作假:" + id)
		}
		codeRunInfo.Header = data.Header
		codeRunInfo.Source = data.Source
		codeRunInfo.Description = ""
	case model.Code_ReRunSubmit:
		data, err := dao.CodeRunGetByIndices(nil, indices)
		if err != nil {
			return codeRunInfo, err
		}
		if data == nil || data.ID != id {
			return codeRunInfo, errors.New("有人想作假:" + id)
		}
		codeRunInfo.Header = data.Header
		codeRunInfo.Source = data.Source
		codeRunInfo.Description = data.Description
	default:
		return codeRunInfo, errors.New("无法解析类别!")
	}
	//
	return codeRunInfo, nil
}

// 代码文件加入待运行队列
func addCodeFile(task CodeRunTaskInfo, retryCount int) error {
	data_byte, err := json.Marshal(task)
	if err != nil {
		return errors.New("加入队列失败:序列化失败!" + err.Error())
	}
	msg := &rocketmq.Message{
		Topic: config.Conf.Code.CodeWaitForRunQueueTopic,
		Body:  data_byte,
	}
	errs := make([]error, 0)
	fn := func() {
		mq.RocketMQProducer.SendAsync(
			context.Background(),
			msg,
			func(ctx context.Context, result []*rocketmq.SendReceipt, err error) {
				// 回调：发送完成后进入
				if err != nil {
					errs = append(errs, err)
					return
				}
			},
		)
	}
	//尝试N次
	for i := 0; i < retryCount; i += 1 {
		fn()
		if len(errs) < i+1 {
			return nil
		}
	}
	return errors.Join(errs...)
}
