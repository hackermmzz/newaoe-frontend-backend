package subrouter

import (
	"newaoe/Src/codeRun/controller"
	"newaoe/Src/filter"

	"github.com/gin-gonic/gin"
)

func CodeRunGroup_Init(group *gin.RouterGroup) {
	coderun_group := group.Group("/coderun")
	coderun_group.Use(filter.FilterCookieCheck())
	//配置路由
	coderun_group.POST("coderun", filter.FilterCodeRun(), controller.CodeRun) //运行代码
	coderun_group.GET("fetchhistory", controller.StudentHistoryGet)           //获取历史记录
	coderun_group.GET("fetchrank", controller.FetchRank)                      //获取排行榜
}
