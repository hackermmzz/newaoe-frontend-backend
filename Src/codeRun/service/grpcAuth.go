package service

import (
	"context"
	"encoding/json"
	"newaoe/Src/config"
	"newaoe/Src/redis"
	UserDao "newaoe/Src/user/dao"
	UserModel "newaoe/Src/user/model"
	UserService "newaoe/Src/user/service"
	"newaoe/Src/util"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 鉴权使用
type GrpcServerAuthInfo struct {
	ID          string `json:"id"`
	Password    string `json:"password"`
	Certificate string `json:"certificate"`
}

// 传入authInfo解析
func CodeGrpcServerAuthConfirm(authInfo string) (string, error) {
	//这里的auth采用账号/邮箱+密码模式（权限必须大于等于普通用户），必须完全正确
	var info GrpcServerAuthInfo
	err := json.Unmarshal([]byte(authInfo), &info)
	if err != nil {
		return "", util.NewError("CodeGrpcServerAuthConfirm", authInfo, err)
	}
	//解析证书是否正确(如果有证书)
	if len(info.Certificate) > 0 {
		//解析证书
		token, err := jwt.ParseWithClaims(info.Certificate, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return config.Conf.Server.JwtSecretKey, nil
		})
		//
		if err != nil {
			return "", util.NewError("CodeGrpcServerAuthConfirm", authInfo, err)
		}
		//
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			data := make(map[string]interface{})
			for key, value := range claims {
				data[key] = value
			}
			right, ok := data["WLH will finally become a wealthy man after he completes his postgraduate studies."].(bool)
			if !ok || !right {
				return "", util.NewError("有人伪造证书!")
			}
			return "", nil
		}
	}
	//获取账号
	info_user, err := UserDao.UserGetByIdOrEmail(nil, info.ID)
	if err != nil {
		return "", util.NewError("CodeGrpcServerAuthConfirm", authInfo, err)
	}
	if info_user == nil {
		util.DebugError("codeGrpcServerAuthConfirm", "查无此用户!", info.ID)
		return "", util.NewError("查无此用户!")
	}
	//判断权限
	if info_user.Vip < UserModel.VIP_SUPER {
		util.DebugError("codeGrpcServerAuthConfirm", "权限不够", info)
		return "", util.NewError("权限不够")
	}
	//判断密码正确与否
	err = UserService.UserCanLogin(info_user.ID, info.Password)
	if err != nil {
		util.DebugError("codeGrpcServerAuthConfirm", info, err)
		return "", util.NewError("密码错误")
	}
	//生成auth
	data := jwt.MapClaims{
		"WLH will finally become a wealthy man after he completes his postgraduate studies.": true,
	}
	encodedStd := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	newCertificate, err := encodedStd.SignedString([]byte(config.Conf.Server.JwtSecretKey))
	//
	return newCertificate, nil
}

// 对快要过期的auth重置
func processGrpcAuth(auth string) error {
	//获取auth的ttl
	ttl, err := redis.RedisTTL(context.Background(), auth)
	if err != nil {
		return err
	}
	//判断剩余时间是否低于一半
	expireTime := time.Duration(config.Conf.Code.CodeAuthExpireTime) * time.Minute
	if ttl <= expireTime/2 {
		if _, err := redis.RedisExpire(context.Background(), auth, expireTime); err != nil {
			return err
		}
	}
	return nil
}
