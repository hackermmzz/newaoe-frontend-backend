package controller

import (
	"net/http"
	"newaoe/Src/common/download/service"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func DownloadAuth(ctx *gin.Context) {
	userInfo := util.GetCtxTookenInfo(ctx)
	if userInfo == nil {
		util.DebugError("请求未携带Cookie!")
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	//
	filePath := strings.TrimLeft(ctx.Param("filepath"), "/")

	//不用管id和vip是否有，如果是公开资源不需要都可以访问
	id, _ := userInfo["id"].(string)
	user, err := dao.UserGet(nil, id)
	if err != nil || user == nil {
		util.DebugError("查询用户失败:", err)
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	vip := user.Vip
	//鉴权
	ok, err := service.DownloadAuth(id, int(vip), filePath)
	if err != nil {
		util.DebugError("DownloadAuth", err)
		util.ResponseNAK_MSG(ctx, "鉴权失败!", nil)
		return
	}
	if !ok {
		util.DebugError("有人尝试越权下载文件:", id, filePath)
		ctx.AbortWithStatus(http.StatusForbidden)
	} else {
		util.ResponseACK_MSG(ctx, "", nil)
	}
}
