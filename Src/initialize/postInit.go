package initialize

import (
	"newaoe/Src/config"
	database "newaoe/Src/databse"
	"newaoe/Src/user/dao"
	"newaoe/Src/user/model"
	"newaoe/Src/util"
)

func PostInit() {
	//所有服务初始化完毕后的预处理
	registerSuperUser()
}

func registerSuperUser() {
	password, err := util.EncodePassword(config.Conf.Other.SuperUserPassword)
	if err != nil {
		util.DebugError("超级用户注册失败:", err.Error())
		return
	}
	for _, userId := range config.Conf.Other.SuperUser {
		if !dao.UserExist(nil, userId) {
			for i := 0; i < 3; i++ {
				session := database.NewSession()
				defer session.Rollback()
				err := session.Begin()
				if err != nil {
					util.DebugError("超级用户", userId, "注册失败", err)
					continue
				}
				info := model.Student{
					Id:         userId,
					Password:   password,
					Email:      config.Conf.Other.SuperUserEmail,
					RegistDate: util.UTC_Time(),
					Avatar:     util.GetRandomAvatar(),
					Vip:        model.VIP_SUPER,
				}
				if !dao.UserAddByStudentInfo(session, info) {
					util.DebugError("超级用户", userId, "注册失败", err)
					continue
				}
				util.DebugSuccess("超级用户:", userId, "注册成功!")
				break
			}
		} else {
			util.Debug("超级用户:", userId, "已经存在!")
		}
	}
	util.DebugSuccess("超级用户注册成毕！")
}
