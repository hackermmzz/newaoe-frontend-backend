package subrouter

import (
	"newaoe/Src/filter"
	"newaoe/Src/vip/controller"

	"github.com/gin-gonic/gin"
)

func VIPGroup_Init(api_group *gin.RouterGroup) {
	//VIP路由
	vipconFire_group := api_group.Group("/vip")
	vipconFire_group.Use(filter.FilterCookieCheck(), filter.FilterVIP())
	{
		routeConfig_VIPConfirm(vipconFire_group)
	}
}

func routeConfig_VIPConfirm(vip_group *gin.RouterGroup) {
	vip_group.GET("fetchStudentInfos", controller.StudentInfosGet)
	vip_group.GET("getstudenthistory", controller.GetStudentHistory)
	vip_group.GET("infoExport", controller.StudentInfoExport)
	vip_group.GET("fetchfeedback", controller.FetchFeedback)
	vip_group.GET("searchstudentbyid", controller.SearchStudentInfo)
	vip_group.POST("resetcommonsubmittime", controller.ResetCommonSubmit)
}
