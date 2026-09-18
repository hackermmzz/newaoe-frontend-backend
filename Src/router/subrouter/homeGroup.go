package subrouter

import (
	"newaoe/Src/filter"
	"newaoe/Src/home/controller"

	"github.com/gin-gonic/gin"
)

func HomeGroup_Init(api_group *gin.RouterGroup) {
	home_group := api_group.Group("/home")
	home_group.Use(filter.FilterCookieCheck())
	{
		routeConfig_Home(home_group)
	}
}

func routeConfig_Home(group *gin.RouterGroup) {
	group.GET("studentInfo", controller.StudentInfoGet)
	group.POST("feedbackUploadAttachment", controller.StudentFeedbckUploadAttachment)
	group.POST("feedbackUpload", controller.StudentFeedbackUpload)
	group.GET("avatarUpdate", controller.AvatarUpdate)
}
