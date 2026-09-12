package Filter

import (
	"newaoe/config"
	"newaoe/util"

	"github.com/gin-gonic/gin"
)

func FilterCookieCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token
		token, err := c.Cookie(config.Conf.Cookie.TokenName)
		if err != nil || token == "" {
			// Cookie 不存在或无效，返回未授权错误
			util.ResponseNAK_MSG(c, "cookie过期或者错误!", "")
			c.Abort()
			return
		}
		//验证cookie是否合法
		if !util.CheckTokenLegal(token) {
			util.ResponseNAK_MSG(c, "cookie过期或者错误!", "")
			c.Abort()
			return
		}
		//
		c.Next()
	}
}
