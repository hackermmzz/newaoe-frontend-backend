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
	vip_group.GET("fetchSubmitRecord", controller.FetchSubmitRecord)
	vip_group.POST("publishAnnouncement", controller.PublishAnnouncement)
	vip_group.GET("rerunAnomalRecord", controller.ReRunAnomalRecord)
	vip_group.GET("blockAllSubmit", controller.BlockAllSubmit)
	vip_group.GET("cancelSubmitBlock", controller.CancelSubmitBlock)
	vip_group.GET("ojVersionUpdate", controller.OJVersionStatusUpdate)
	vip_group.GET("teacherAdd", controller.TeacherAdd)
	vip_group.GET("runAllAssessmentSubmit", controller.RunAllAssessmentSubmit)
	vip_group.GET("getassessmenthistory", controller.GetAssessmentHistory)
	vip_group.GET("fetchassessmentrank", controller.FetchAssessmentRank)
	vip_group.GET("uploadLatestAOEPackage", controller.UploadLatestAOEPackage)
}
