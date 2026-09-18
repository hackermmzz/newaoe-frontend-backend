package router

import (
	"newaoe/Src/router/subrouter"

	"github.com/gin-gonic/gin"
)

func RouterConfig(Engine *gin.Engine) {
	api_gorup := Engine.Group("/api")
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

//	func routeConfig_UploadConfirm(uploadconFirm_group *gin.RouterGroup) {
//		uploadconFirm_group.POST("/:category", Upload.FileUploadConfirm) //告诉后端上传好了
//	}
// func routeConfig_Code(code_group *gin.RouterGroup) {
// 	code_group.POST("/CodeReRun", Filter.FilterCookieCheck(), Filter.FilterCodeRun(), Filter.FilterLimitCodeUploadOrRun(), Code.CodeReRun) //这个要使用cookie检测中间件
// 	//
// 	/*这里使用grpc代替之前的http
// 	code_group.GET("/CodeGet", Filter.FilterCodeRunServerCheck(), Code.CodeGetService)                      //获取代码
// 	code_group.POST("/CodeRunStatusPost", Filter.FilterCodeRunServerCheck(), Code.CodeRunStatusPostService) //处理代码运行状态上传
// 	*/
// }
