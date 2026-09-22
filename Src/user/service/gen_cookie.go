package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"newaoe/Src/config"
	"newaoe/Src/user/dao"
	"newaoe/Src/util"
	"time"
)

func UserGenCookie(id string, ip string) (string, error) {
	//获取用户信息
	user := dao.UserGet(nil, id)
	if user == nil {
		return "", errors.New("数据库异常!")
	}
	//
	data, err := json.Marshal(*user)
	if err != nil {
		return "", errors.New("UserSetCookie marshal err" + err.Error())
	}
	var info map[string]interface{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return "", errors.New("UserSetCookie unmarshal err" + err.Error())
	}
	tooken_value := util.GenerateCookieToken(ip,
		config.Conf.Cookie.ExpireTime,
		config.Conf.Server.JwtSecretKey,
		info,
	)
	//
	if tooken_value == "" {
		return "", errors.New("生成token失败!")
	}
	//
	cookie := &http.Cookie{
		Name:   config.Conf.Cookie.TokenName,
		Value:  tooken_value,
		MaxAge: config.Conf.Cookie.ExpireTime,
		Path:   "/",
		//Domain:   "example.com",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  util.UTC_Time().Add(time.Duration(config.Conf.Cookie.ExpireTime) * time.Second),
	}
	cookie_str := cookie.String()
	return cookie_str, nil
}
