package Upload

import (
	"context"
	"encoding/json"
	"fmt"
	"newaoe/config"
	"newaoe/dao"
	"newaoe/util"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func CodeUpload(ctx *gin.Context, userInfo map[string]string) {

	id, _ := userInfo["id"]
	//限制提交次数
	if limitCodeUpload(id) {
		util.ResponseNAK_MSG(ctx, "当天3次提交次数已消耗尽!", nil)
		return
	}
	//
	prefix := util.UTC_Time().Format("2006_01_02_150405000")
	base_dir := path.Join(config.Conf.OSS.PrivateBaseFolder, id, config.Conf.User.UserCodeFolder, prefix)
	//
	header := path.Join(base_dir, "mmzz.h")
	source := path.Join(base_dir, "mmzz.cpp")
	description := path.Join(base_dir, "mmzz.txt")
	duration := time.Duration(config.Conf.Other.CodeCommonUploadUrlExpireTime) * time.Second
	UploadFile(ctx, []string{header, source, description}, []time.Duration{duration, duration, duration})
}

func limitCodeUpload(id string) bool {
	//
	key := fmt.Sprintf("CommonUploadTimes_%v", id)
	time, exist := dao.RedisGet(context.Background(), key)
	curSubmitTime := 0
	if exist {
		err := json.Unmarshal(time, &curSubmitTime)
		if err != nil {
			util.Debug("limitCodeUpload Unmarshal Error!", err)
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
