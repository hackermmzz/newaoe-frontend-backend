package data

import (
	"newaoe/dao"

	"xorm.io/xorm"
)

/*
这个我没加到redis是因为我使用了jwt技术
他是无状态的
*/

// 判断用户是否允许注册(不允许未授权用户注册)
func UserIdLegal(id string) bool {
	return true
}

// 添加用户
func UserAdd(session *xorm.Session, id string, password string, email string) bool {
	ok := dao.UserAdd(session, id, password, email)
	return ok
}

// 更新用户数据
func UserUpdate(session *xorm.Session, id string, user dao.Student) bool {
	ok := dao.UserUpdate(session, id, user)
	return ok
}

// 判断用户是否已经注册
func UserExist(id string) bool {
	return dao.UserExist(id)
}

// 判断邮箱是否已经使用
func EmailExist(email string) bool {
	return dao.EmailExist(email)
}

// 获取指定用户的数据
func UserGet(id string) *dao.Student {
	return dao.UserGet(id)
}

// 获取用户密码(编码后的密码)
func UserGetPassword(id string) string {
	data := UserGet(id)
	if data == nil {
		return ""
	}
	return data.Password
}

func UserGetEmail(id string) string {
	data := UserGet(id)
	if data == nil {
		return ""
	}
	return data.Email
}

func UserGetAvatar(id string) string {
	data := UserGet(id)
	if data == nil {
		return ""
	}
	return data.Avatar
}

func UserResetPassword(session *xorm.Session, id string, new_password string) bool {
	return UserUpdate(session, id, dao.Student{
		Password: new_password,
	})
}

func UserResetAvatar(session *xorm.Session, id string, avatarFile string) bool {
	return UserUpdate(session, id, dao.Student{
		Avatar: avatarFile,
	})
}
