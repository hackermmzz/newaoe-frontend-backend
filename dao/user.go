package dao

import (
	"newaoe/util"
	"time"

	"xorm.io/xorm"
)

type Student struct {
	Id         string    `xorm:"id"`
	Password   string    `xorm:"password"`
	Email      string    `xorm:"email"`
	Avatar     string    `xorm:"avatar"`
	RegistDate time.Time `xorm:"registDate"`
}

func (Student) TableName() string {
	return "Student" // 这里返回你想要的表名
}

// 判断用户是否已经注册
func UserExist(id string) bool {
	has, err := DB.Where("id = ?", id).Exist(&Student{})
	if (!has) || (err != nil) {
		if err != nil {
			util.Debug("UserExist: ", err)
		}
		return false
	}
	return true
}

// 判断邮箱是否已经使用
func EmailExist(email string) bool {
	has, err := DB.Where("email = ?", email).Exist(&Student{})
	if (!has) || (err != nil) {
		if err != nil {
			util.Debug("EmailExist: ", err)
		}
		return false
	}
	return true
}

// 添加用户
func UserAdd(session *xorm.Session, id string, password string, email string) bool {
	data := &Student{
		Id:         id,
		Password:   password,
		Email:      email,
		RegistDate: util.UTC_Time(),
		Avatar:     util.GetRandomAvatar(),
	}
	_, err := session.Insert(data)
	if err != nil {
		util.Debug("AddUser:", err)
		return false
	}
	return true
}

// 更新用户数据
func UserUpdate(session *xorm.Session, id string, user Student) bool {
	_, err := session.Where("id = ?", id).Update(&user)
	if err != nil {
		util.Debug("UserUpdate:", err)
		return false
	}
	return true
}

// 获取指定用户的数据
func UserGet(id string) *Student {
	var user Student
	has, err := DB.Where("id = ? ", id).Get(&user)
	if err != nil || !has {
		util.Debug("UserGet:", err)
		return nil
	}
	return &user
}
