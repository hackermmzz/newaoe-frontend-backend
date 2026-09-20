package subrouter

import (
	"newaoe/Src/user/controller"

	"github.com/gin-gonic/gin"
)

func UserGroup_Init(api_group *gin.RouterGroup) {
	//用户路由组(登陆、注册、修改密码等)
	user_group := api_group.Group("/user")
	{
		routeConfig_User(user_group)
	}
}

func routeConfig_User(group *gin.RouterGroup) {
	//登陆
	group.POST("login", controller.UserLogin)
	//普通用户注册
	group.POST("regist", controller.UserRegist)
	//游客注册
	group.POST("touristregist", controller.TouristUserRegist)
	//发送注册验证码
	group.POST("registCode", controller.UserRegistCodeSend)
	//重置密码
	group.POST("resetPassword", controller.UserResetPassword)
	//发送重置密码验证码
	group.POST("resetPasswordCode", controller.UserPasswordForgetCodeSend)
}
