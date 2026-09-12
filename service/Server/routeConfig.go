package Server

import (
	"fmt"
	"newaoe/config"
	"newaoe/service/Code"
	"newaoe/service/Download"
	"newaoe/service/Home"
	rank "newaoe/service/Rank"
	"newaoe/service/Server/Filter"
	"newaoe/service/Upload"
	"newaoe/service/User"

	"github.com/gin-gonic/gin"
)

func routeConfig() {
	api_gorup := Engine.Group("/api")
	{
		//用户路由组(登陆、注册、修改密码等)
		user_group := api_gorup.Group("/user")
		{
			routeConfig_User(user_group)
		}
		//界面路由组
		home_group := api_gorup.Group("/home")
		home_group.Use(Filter.FilterCookieCheck())
		{
			routeConfig_Home(home_group)
		}
		//排行榜路由
		rank_group := api_gorup.Group("/rank")
		rank_group.Use(Filter.FilterCookieCheck())
		{
			routeConfig_Rank(rank_group)
		}
		//下载路由组
		download_group := api_gorup.Group("/download")
		download_group.Use(Filter.FilterCookieCheck())
		{
			routeConfig_Download(download_group)
		}
		//上传路由组
		upload_group := api_gorup.Group("/upload")
		upload_group.Use(Filter.FilterCookieCheck())
		{
			routeConfig_Upload(upload_group)
		}
		//上传确认路由
		uploadconFirm_group := api_gorup.Group("/uploadconfirm")
		uploadconFirm_group.Use(Filter.FilterCookieCheck())
		{
			routeConfig_UploadConfirm(uploadconFirm_group)
		}
		//接口路由
		code_group := api_gorup.Group("/code")
		{
			routeConfig_Code(code_group)
		}
	}
}

func routeConfig_Rank(rank_group *gin.RouterGroup) {
	rank_group.GET("/fetchrank", rank.FetchRank)
}

func routeConfig_UploadConfirm(uploadconFirm_group *gin.RouterGroup) {
	uploadconFirm_group.POST("/:category", Upload.FileUploadConfirm) //告诉后端上传好了
}
func routeConfig_Code(code_group *gin.RouterGroup) {
	code_group.POST("/CodeReRun", Filter.FilterCookieCheck(), Filter.FilterCodeRun(), Filter.FilterLimitCodeUploadOrRun(), Code.CodeReRun) //这个要使用cookie检测中间件
	//
	/*这里使用grpc代替之前的http
	code_group.GET("/CodeGet", Filter.FilterCodeRunServerCheck(), Code.CodeGetService)                      //获取代码
	code_group.POST("/CodeRunStatusPost", Filter.FilterCodeRunServerCheck(), Code.CodeRunStatusPostService) //处理代码运行状态上传
	*/
}

func routeConfig_Upload(upload_group *gin.RouterGroup) {
	//普通文件的上传
	upload_group.POST(fmt.Sprintf("/%v", config.Conf.User.UserCodeFolder),
		Filter.FilterCodeRun(), Filter.FilterLimitCodeUploadOrRun(), Upload.FileUpload)
	//考核代码文件上传
	upload_group.POST(fmt.Sprintf("/%v", config.Conf.Other.AssessmentCodeUploadCategory), Filter.FilterCodeRun(), Filter.FilterAssessmentSubmission(), Upload.FileUpload)
	//头像提交
	upload_group.POST(fmt.Sprintf("/%v", config.Conf.User.UserAvatarFolder), Upload.FileUpload)
}

func routeConfig_Download(download_group *gin.RouterGroup) {
	download_group.GET("/:category/*filepath", Download.FileDownload)
}

func routeConfig_User(group *gin.RouterGroup) {
	//登陆
	group.POST("/login", User.UserLogin)
	//注册
	group.POST("/regist", User.UserRegist)
	//发送注册验证码
	group.POST("/registCode", User.UserRegistCodeSend)
	//重置密码
	group.POST("/resetPassword", User.UserResetPassword)
	//发送重置密码验证码
	group.POST("/resetPasswordCode", User.UserPasswordForgetCodeSend)
}

func routeConfig_Home(group *gin.RouterGroup) {
	group.GET("studentInfo", Home.StudentInfoGet)
	group.POST("feedback", Home.StudentFeedback)
	group.GET("gethistory", Home.StudentHistoryGet)
	group.GET("getteacher", Home.StudentGetTeacher)
}
