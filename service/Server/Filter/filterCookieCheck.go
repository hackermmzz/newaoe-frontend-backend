package Filter

import (
	"encoding/json"
	"newaoe/config"
	"newaoe/dao"
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
		//判断它是否能解析出完整的Student信息
		var student dao.Student
		data := util.GetTokenInfo(token)
		dt0, _ := json.Marshal(data)
		err = json.Unmarshal(dt0, &student)
		if err != nil || student.Id == "" {
			util.DebugError("有人伪造Token信息!", data)
			util.ResponseNAK_MSG(c, "cookie过期或者错误!", "")
			c.Abort()
			return
		}
		//验证cookie是否合法
		expire_time, ok := data["expire_time"].(string)
		if !ok || !util.CheckTokenLegal(expire_time) {
			util.ResponseNAK_MSG(c, "cookie过期或者错误!", "")
			c.Abort()
			return
		}
		//
		c.Next()
	}
}
