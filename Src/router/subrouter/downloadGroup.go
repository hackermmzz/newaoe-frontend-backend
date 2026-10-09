package subrouter

import (
	"fmt"
	"newaoe/Src/common/download/controller"
	"newaoe/Src/config"
	"newaoe/Src/filter"

	"github.com/gin-gonic/gin"
)

func DownloadGroup_Init(api_group *gin.RouterGroup) {
	g := api_group.Group("/download")
	g.Use(filter.FilterCookieCheck()) //检测cookie
	//配置下载路由
	g.GET("/*filepath", controller.FileDownload)

	//鉴权路由(无需鉴定cookie，无则不管，铭感数据解析失败自然不会给他下载)
	gg := api_group.Group("downloadauth")
	gg.GET(fmt.Sprintf("%s/*filepath", config.Conf.OSS.BucketName), controller.DownloadAuth)
}
