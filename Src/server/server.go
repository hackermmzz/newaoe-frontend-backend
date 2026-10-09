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
	gin.SetMode(config.Conf.Server.ServeMode)
	Engine = gin.Default()
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
