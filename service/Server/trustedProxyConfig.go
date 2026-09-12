package Server

import "newaoe/config"

func trustedProxyConfig() {
	err := Engine.SetTrustedProxies(config.Conf.Server.TrustedProxy)
	if err != nil {
		panic("配置trustedProxy出现问题:" + err.Error())
	}

}
