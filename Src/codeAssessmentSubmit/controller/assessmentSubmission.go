package controller

import (
	"context"
	"errors"
	"fmt"
	"newaoe/Src/codeAssessmentSubmit/service"
	"newaoe/Src/redis"
	"newaoe/Src/util"
	"time"

	"github.com/gin-gonic/gin"
)

func AssessmentSubmissionSubmit(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	id := userInfo["id"].(string)
	//获取上传链接以及信息key
	urls, key, err := service.AssessmentSubmit(id)
	//回复
	if err != nil {
		util.DebugError("AssessmentSubmissionSubmit失败!", err)
		util.ResponseNAK_MSG(ctx, "服务器异常", nil)
		return
	}
	//生成pending的key，防止客户端恶意重复请求生成key
	pendingKey := fmt.Sprintf("AssessmentSubmissionPending_%s", id)
	ok, err := redis.RedisSet(context.Background(), pendingKey, key, time.Duration(3)*time.Minute)
	if !ok || err != nil {
		if err == nil {
			err = errors.New("生成pending key失败!")
		}
		util.DebugError("AssessmentSubmissionSubmit失败!", err)
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//
	util.ResponseACK_MSG(ctx, "获取链接成功!",
		map[string]interface{}{
			"urls": urls,
			"key":  key,
		})
}
