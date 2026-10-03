package router

import (
	"newaoe/Src/filter"
	"newaoe/Src/router/subrouter"

	"github.com/gin-gonic/gin"
)

func RouterConfig(Engine *gin.Engine) {
	api_gorup := Engine.Group("/api")
	configApiRouter(api_gorup)
	{
		//用户路由组
		subrouter.UserGroup_Init(api_gorup)
		//Home界面路由组
		subrouter.HomeGroup_Init(api_gorup)
		//代码运行路由
		subrouter.CodeRunGroup_Init(api_gorup)
		//代码提交路由
		subrouter.CodeSubmitGroup_Init(api_gorup)
		//VIP路由
		subrouter.VIPGroup_Init(api_gorup)
		//下载路由
		subrouter.DownloadGroup_Init(api_gorup)
	}
}

func configApiRouter(group *gin.RouterGroup) {
	group.GET("checkLoginStatus",filter.FilterCookieCheck())
}
