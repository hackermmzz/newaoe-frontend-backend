package server

import (
	"newaoe/Src/config"

	"github.com/gin-gonic/gin"
)

func trustedProxyConfig(Engine *gin.Engine) {
	err := Engine.SetTrustedProxies(config.Conf.Server.TrustedProxy)
	if err != nil {
		panic("配置trustedProxy出现问题:" + err.Error())
	}

}
