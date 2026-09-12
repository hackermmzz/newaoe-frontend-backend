package User

import (
	"net/http"
	"newaoe/config"
	data "newaoe/service/Data"
	"newaoe/util"
	"time"

	"github.com/gin-gonic/gin"
)

func UserSetCookie(ctx *gin.Context, id string) string {
	//获取用户信息
	user := data.UserGet(id)
	if user == nil {
		util.Debug(ctx, "数据库异常!", "")
		return ""
	}
	//
	tooken_value := util.GenerateCookieToken(ctx, id, user.Email, user.RegistDate, config.Conf.Cookie.ExpireTime, config.Conf.Server.JwtSecretKey)
	//
	if tooken_value == "" {
		util.Debug(ctx, "生成token失败!", "")
		return ""
	}
	//
	cookie := &http.Cookie{
		Name:   config.Conf.Cookie.TokenName,
		Value:  tooken_value,
		MaxAge: config.Conf.Cookie.ExpireTime,
		Path:   "/api",
		//Domain:   "example.com",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  util.UTC_Time().Add(time.Duration(config.Conf.Cookie.ExpireTime) * time.Second),
	}
	cookie_str := cookie.String()
	ctx.Header("set-cookie", cookie_str)
	return tooken_value
}
