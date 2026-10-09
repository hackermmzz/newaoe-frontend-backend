package server

import (
	"fmt"
	"newaoe/Src/config"
	"newaoe/Src/router"
	"newaoe/Src/util"

	"github.com/gin-gonic/gin"
)

var (
	Engine *gin.Engine
)

func ServerInit() {
	Engine = gin.Default()
	Engine.Use(gin.Recovery()) // 关键：添加 Recovery 捕获 panic 并输出日志
	Engine.Use(gin.Logger())   // 输出请求和错误日志
	gin.SetMode(config.Conf.Server.ServeMode)
	//
	// 配置CORS，处理OPTIONS预检请求
	cors_Config(Engine)
	//配置路由以及配置过滤器
	router.RouterConfig(Engine)
	//受信任代理配置
	trustedProxyConfig(Engine)
	//
	if err := Engine.Run(fmt.Sprintf(":%d", config.Conf.Server.ServerListenPort)); err != nil {
		util.DebugError("服务器启动失败: " + err.Error())
	}
	util.DebugSuccess("服务器启动成功")
}
