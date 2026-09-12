package Server

import (
	"fmt"
	"newaoe/config"
	"newaoe/util"

	"github.com/gin-gonic/gin"
)

var (
	Engine *gin.Engine
)

func ServeInit() {
	Engine = gin.Default()
	Engine.Use(gin.Recovery()) // 关键：添加 Recovery 捕获 panic 并输出日志
	Engine.Use(gin.Logger())   // 输出请求和错误日志
	gin.SetMode(config.Conf.Server.ServeMode)
	//
	// 配置CORS，处理OPTIONS预检请求
	cors_Config()
	//配置路由以及配置过滤器
	routeConfig()
	//受信任代理配置
	trustedProxyConfig()
	//
	if err := Engine.Run(fmt.Sprintf(":%d", config.Conf.Server.ServerListenPort)); err != nil {
		util.Debug("服务器启动失败: " + err.Error())
	}
	util.Debug("服务器启动成功")
}
