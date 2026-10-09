package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/util"
	"time"
)

func GetOneCodeTask() (*CodeRunTaskInfo, error) {
	//这里要做幂等，防止这个消息已经被消费过了
	errs := make([]error, 0)
	for i := 0; i < 10; i += 1 {
		ret, err := pullOneCodefunc()
		if err == nil {
			return ret, nil
		}
		if ret == nil {
			continue
		}
		errs = append(errs, err)
	}
	//
	return nil, errors.Join(errs...)
}

func pullOneCodefunc() (*CodeRunTaskInfo, error) {
	var codeinfo CodeRunTaskInfo
	//拉取数据
	msgs, err := CodeWaitForRunPullConsumer.Receive(
		context.Background(),
		1,             // 最多拉1条
		5*time.Minute, // 消息不可见时间
	)
	if err != nil {
		return nil, err
	}
	//ACK
	msg := msgs[0]
	err = CodeWaitForRunPullConsumer.Ack(context.Background(), msg)
	if err != nil {
		return nil, err
	}
	//发序列化数据
	err = json.Unmarshal(msg.GetBody(), &codeinfo)
	if err != nil {
		return nil, util.NewError(err)
	}
	//判断是否已经处理过了(只有CodeRunning表里面没有，且CodeRun表有才算成功跑结束，取反就是下面这个)
	has0, err := dao.CodeRunningExist(nil, codeinfo.Indices)
	has1 := dao.CodeRunExist(nil, codeinfo.Indices)
	if err != nil {
		return nil, util.NewError(err)
	}
	if has0 || !has1 {
		//
		util.DebugSuccess(fmt.Sprintf("recv msgId=%s receipt=%s. OJ successfully get one code! ID:%v, Indices:%v",
			msg.GetMessageId(),
			msg.GetReceiptHandle(),
			codeinfo.ID,
			codeinfo.Indices,
		))
		return &codeinfo, nil
	}
	return nil, nil
}
