package subrouter

import (
	"newaoe/Src/common/download/controller"
	"newaoe/Src/filter"

	"github.com/gin-gonic/gin"
)

func DownloadGroup_Init(api_group *gin.RouterGroup) {
	g := api_group.Group("/download")
	g.Use(filter.FilterCookieCheck()) //检测cookie
	//配置下载路由
	g.GET("/*filepath", controller.FileDownload)
}
