package User

import (
	"encoding/json"
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
		util.DebugError(ctx, "数据库异常!", "")
		return ""
	}
	//
	data, err := json.Marshal(*user)
	if err != nil {
		util.Debug("UserSetCookie marshal err", err)
		return ""
	}
	var info map[string]interface{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		util.Debug("UserSetCookie unmarshal err", err)
		return ""
	}
	tooken_value := util.GenerateCookieToken(ctx,
		config.Conf.Cookie.ExpireTime,
		config.Conf.Server.JwtSecretKey,
		info,
	)
	//
	if tooken_value == "" {
		util.DebugError(ctx, "生成token失败!", "")
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
