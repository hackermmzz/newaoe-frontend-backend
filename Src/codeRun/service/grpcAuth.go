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
)

// 鉴权使用
type GrpcServerAuthInfo struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Auth     string `json:"auth"`
}

func CodeGrpcServerAuthConfirm(auth string) (bool, string) {
	//这里的auth采用账号/邮箱+密码模式（权限必须大于等于普通用户），必须完全正确
	var info GrpcServerAuthInfo
	err := json.Unmarshal([]byte(auth), &info)
	if err != nil {
		util.DebugError("codeGrpcServerAuthConfirm", auth, err)
		return false, ""
	}
	//获取auth判断是否可信
	if redis.RedisExist(context.Background(), info.Auth) {
		//获取过期时间
		if err := processGrpcAuth(info.Auth); err != nil {
			util.DebugError("CodeGrpcServerAuthConfirm processGrpcAuth", err)
		}
		//
		return true, info.Auth
	}
	//获取账号
	info_user, err := UserDao.UserGetByIdOrEmail(nil, info.ID)
	if err != nil {
		util.DebugError("UserGetByIdOrEmail:", err)
	}
	if info_user == nil {
		util.DebugError("codeGrpcServerAuthConfirm", "查无此用户!", info.ID)
		return false, ""
	}
	//判断权限
	if info_user.Vip < UserModel.VIP_SUPER {
		util.DebugError("codeGrpcServerAuthConfirm", "权限不够", info)
		return false, ""
	}
	//判断密码正确与否
	err = UserService.UserCanLogin(info_user.Id, info.Password)
	if err != nil {
		util.DebugError("codeGrpcServerAuthConfirm", info, err)
		return false, ""
	}
	//加入redis
	newAuth := util.UUID()
	redis.RedisSet(context.Background(), newAuth, nil, time.Duration(config.Conf.Code.CodeAuthExpireTime)*time.Minute)
	//
	return true, newAuth
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
		if !redis.RedisExpire(context.Background(), auth, expireTime) {
			return util.NewError("重置auth的ttl失败!")
		}
	}
	return nil
}
