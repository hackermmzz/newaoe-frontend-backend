package controller

import (
	"fmt"
	"newaoe/Src/common/upload"
	"newaoe/Src/config"
	"newaoe/Src/map/dao"
	"newaoe/Src/util"
	"newaoe/Src/vip/service"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadLatestAOEPackage(ctx *gin.Context) {
	detail := ctx.Query("detail")
	//直接操作oss
	filepath := filepath.Join(config.Conf.OSS.PublicBaseFolder, config.Conf.Other.PublicOtherFolder, fmt.Sprintf("newaoe-lhw-%v.zip", util.UUID()))
	url := upload.UploadFile(filepath, time.Duration(1)*time.Hour)
	if url == "" {
		util.ResponseNAK_MSG(ctx, "生成上传链接失败!", nil)
		return
	}
	//更新数据库
	err := dao.MapUpdateOrInsert(nil, "newaoe-latest-version", filepath)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//发布公告
	err = service.PublishAnnouncement("NewAOE新版本推送通知", detail, true)
	if err != nil {
		util.ResponseNAK_MSG(ctx, err.Error(), nil)
		return
	}
	//
	util.ResponseACK_MSG(ctx, "生成链接成功!", map[string]interface{}{"url": url})
}
