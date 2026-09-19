package dao

import (
	database "newaoe/Src/databse"
	"newaoe/Src/user/model"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func UserResetAvatar(session *xorm.Session, id string, avatarFile string) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	return UserUpdate(session, id, model.Student{
		Avatar: avatarFile,
	})
}

// 判断用户是否已经注册
func UserExist(session *xorm.Session, id string) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	has, err := session.Where("id = ?", id).Exist(&model.Student{})
	if (!has) || (err != nil) {
		if err != nil {
			util.DebugError("UserExist: ", err)
		}
		return false
	}
	return true
}

// 判断邮箱是否已经使用
func EmailExist(session *xorm.Session, email string) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	has, err := session.Where("email = ?", email).Exist(&model.Student{})
	if (!has) || (err != nil) {
		if err != nil {
			util.DebugError("EmailExist: ", err)
		}
		return false
	}
	return true
}

// 添加用户
func UserAdd(session *xorm.Session, id string, password string, email string) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data := &model.Student{
		Id:         id,
		Password:   password,
		Email:      email,
		RegistDate: util.UTC_Time(),
		Avatar:     util.GetRandomAvatar(),
	}
	_, err := session.Insert(data)
	if err != nil {
		util.DebugError("AddUser:", err)
		return false
	}
	return true
}

// 添加用户
func UserAddByStudentInfo(session *xorm.Session, info model.Student) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	_, err := session.Insert(info)
	if err != nil {
		util.DebugError("AddUser:", err)
		return false
	}
	return true
}

// 更新用户数据
func UserUpdate(session *xorm.Session, id string, user model.Student) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	_, err := session.Where("id = ?", id).Update(&user)
	if err != nil {
		util.DebugError("UserUpdate:", err)
		return false
	}
	return true
}

// 获取指定用户的数据
func UserGet(session *xorm.Session, id string) *model.Student {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var user model.Student
	has, err := session.Where("id = ? ", id).Get(&user)
	if err != nil || !has {
		util.DebugError("UserGet:", err)
		return nil
	}
	return &user
}

// 根据注册时间升序排序（不包括end)
func UserGetByRangeOrderByRegistData(session *xorm.Session, beg int, end int) []model.Student {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if beg < 0 || beg >= end {
		return nil
	}
	var list []model.Student
	// 按regist_date升序
	// LIMIT beg, end-beg  等价于 offset beg limit count
	err := session.Asc("registDate").Limit(end-beg, beg).Find(&list)
	if err != nil {
		util.DebugError("UserGetByRangeOrderByRegistData error:", err)
		return nil
	}
	return list
}

// 获取指定多个用户的数据
func UserGetByIDs(session *xorm.Session, ids []string) []model.Student {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//

	var users []model.Student

	if len(ids) == 0 {
		return users
	}

	err := session.In("id", ids).Find(&users)
	if err != nil {
		util.DebugError("UserGetByIDs:", err)
		return nil
	}

	return users
}

// 获取所有学生数据(谨慎使用)
func UserGetAll(session *xorm.Session) []model.Student {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var list []model.Student
	// Find 直接把查询结果填充到切片
	err := session.Find(&list)
	if err != nil {
		// 出错返回空
		util.DebugError("UserGetAll:", err)
		return nil
	}
	return list
}

// 判断用户是否允许注册(不允许未授权用户注册)
func UserIdLegal(id string) bool {
	return true
}

// 获取用户密码(编码后的密码)
func UserGetPassword(session *xorm.Session, id string) string {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data := UserGet(session, id)
	if data == nil {
		return ""
	}
	return data.Password
}

func UserGetEmail(session *xorm.Session, id string) string {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data := UserGet(session, id)
	if data == nil {
		return ""
	}
	return data.Email
}

func UserGetAvatar(session *xorm.Session, id string) string {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data := UserGet(session, id)
	if data == nil {
		return ""
	}
	return data.Avatar
}

func UserResetPassword(session *xorm.Session, id string, new_password string) bool {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	return UserUpdate(session, id, model.Student{
		Password: new_password,
	})
}
