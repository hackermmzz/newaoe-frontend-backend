package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"newaoe/Src/codeRun/dao"
	"newaoe/Src/oss"
	"newaoe/Src/util"
	"time"
)

func GetOneCodeTask() (*CodeRunTaskInfo, []string, error) {
	//这里要做幂等，防止这个消息已经被消费过了
	errs := make([]error, 0)
	for i := 0; i < 10; i += 1 {
		info, urls, err := pullOneCodefunc()
		if err == nil {
			return info, urls, nil
		}
		errs = append(errs, err)
	}
	//
	return nil, nil, errors.Join(errs...)
}

func pullOneCodefunc() (info *CodeRunTaskInfo, urls []string, err error) {
	var codeinfo CodeRunTaskInfo
	//拉取数据
	msgs, err := CodeWaitForRunPullConsumer.Receive(
		context.Background(),
		1,             // 最多拉1条
		5*time.Minute, // 消息不可见时间
	)
	if err != nil {
		return nil, nil, err
	}
	if len(msgs) == 0 {
		return nil, nil, nil
	}
	//ACK
	msg := msgs[0]
	defer func() {
		if err == nil {
			if err := CodeWaitForRunPullConsumer.Ack(context.Background(), msg); err != nil {
				util.DebugError("pullOneCodefunc:ACK失败!", err)
			}
		}
	}()
	//发序列化数据
	err = json.Unmarshal(msg.GetBody(), &codeinfo)
	if err != nil {
		return nil, nil, util.NewError(err)
	}
	//判断是否已经处理过了(只有CodeRunning表里面没有，且CodeRun表有才算成功跑结束，取反就是下面这个)
	has0, err := dao.CodeRunningExist(nil, codeinfo.Indices)
	if err != nil {
		return nil, nil, util.NewError(err)
	}
	has1, err := dao.CodeRunExist(nil, codeinfo.Indices)
	if err != nil {
		return nil, nil, util.NewError(err)
	}
	if err != nil {
		return nil, nil, util.NewError(err)
	}
	if has0 || !has1 {
		//
		util.DebugSuccess(fmt.Sprintf("recv msgId=%s receipt=%s. OJ successfully get one code! ID:%v, Indices:%v",
			msg.GetMessageId(),
			msg.GetReceiptHandle(),
			codeinfo.ID,
			codeinfo.Indices,
		))
		//获取下载链接
		source_dir := codeinfo.Source
		header_dir := codeinfo.Header
		//获取链接（这里不走路由了，直接后端获取链接）
		duration := time.Duration(30) * time.Minute
		urlGet := []string{}
		for i := 0; i < 3 && len(urlGet) == 0; i += 1 {
			urlGet = oss.OssGetDownloadFileUrls([]string{header_dir, source_dir}, duration, false)
		}
		if len(urlGet) == 0 {
			return nil, nil, util.NewError("生成下载链接失败!")
		}
		//
		return &codeinfo, urlGet, nil
	}
	return nil, nil, nil
}
