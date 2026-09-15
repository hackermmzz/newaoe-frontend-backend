package Service

import (
	"newaoe/config"
	"newaoe/dao"
	"newaoe/service/Code"
	data "newaoe/service/Data"
	"newaoe/service/Email"
	"newaoe/service/Home"
	"newaoe/util"
)

func ServiceInit() {
	//初始化邮箱连接
	Email.EmailSenderInit()
	//初始化反馈推送
	Home.FeedbackInit()
	//初始化代码运行状态服务
	Code.CodeRunServiceInit()
}

func PostProcess() {
	//注册所有超级用户
	registerSuperUser()
}

func registerSuperUser() {
	password, err := util.EncodePassword(config.Conf.Other.SuperUserPassword)
	if err != nil {
		util.DebugError("超级用户注册失败:", err.Error())
		return
	}
	for _, userId := range config.Conf.Other.SuperUser {
		if !data.UserExist(userId) {
			for i := 0; i < 3; i++ {
				session := dao.DB.NewSession()
				defer session.Rollback()
				err := session.Begin()
				if err != nil {
					util.DebugError("超级用户", userId, "注册失败", err)
					continue
				}
				info := dao.Student{
					Id:         userId,
					Password:   password,
					Email:      config.Conf.Other.SuperUserEmail,
					RegistDate: util.UTC_Time(),
					Avatar:     util.GetRandomAvatar(),
					Vip:        dao.VIP_SUPER,
				}
				if !data.UserAddByStudentInfo(session, info) {
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
