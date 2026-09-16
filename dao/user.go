package dao

import (
	"newaoe/util"
	"time"

	"xorm.io/xorm"
)

// VIP等级
var (
	VIP_NONE  = 0 //没有特权
	VIP_SUPER = 1 //特权用户
)

type Student struct {
	Id         string    `json:"id" xorm:"id"`
	Password   string    `json:"password" xorm:"password"`
	Email      string    `json:"email" xorm:"email"`
	Avatar     string    `json:"avatar" xorm:"avatar"`
	RegistDate time.Time `json:"registDate" xorm:"registDate"`
	Vip        int       `json:"vip" xorm:"vip"`
}

func (Student) TableName() string {
	return "Student" // 这里返回你想要的表名
}

// 判断用户是否已经注册
func UserExist(id string) bool {
	has, err := DB.Where("id = ?", id).Exist(&Student{})
	if (!has) || (err != nil) {
		if err != nil {
			util.DebugError("UserExist: ", err)
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
			util.DebugError("EmailExist: ", err)
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
		util.DebugError("AddUser:", err)
		return false
	}
	return true
}

// 添加用户
func UserAddByStudentInfo(session *xorm.Session, info Student) bool {
	_, err := session.Insert(info)
	if err != nil {
		util.DebugError("AddUser:", err)
		return false
	}
	return true
}

// 更新用户数据
func UserUpdate(session *xorm.Session, id string, user Student) bool {
	_, err := session.Where("id = ?", id).Update(&user)
	if err != nil {
		util.DebugError("UserUpdate:", err)
		return false
	}
	return true
}

// 获取指定用户的数据
func UserGet(id string) *Student {
	var user Student
	has, err := DB.Where("id = ? ", id).Get(&user)
	if err != nil || !has {
		util.DebugError("UserGet:", err)
		return nil
	}
	return &user
}

// 根据注册时间升序排序（不包括end)
func UserGetByRangeOrderByRegistData(beg int, end int) []Student {
	if beg < 0 || beg >= end {
		return nil
	}
	var list []Student
	// 按regist_date升序
	// LIMIT beg, end-beg  等价于 offset beg limit count
	err := DB.Asc("registDate").Limit(end-beg, beg).Find(&list)
	if err != nil {
		util.DebugError("UserGetByRangeOrderByRegistData error:", err)
		return nil
	}
	return list
}

// 获取指定多个用户的数据
func UserGetByIDs(ids []string) []Student {

	var users []Student

	if len(ids) == 0 {
		return users
	}

	err := DB.In("id", ids).Find(&users)
	if err != nil {
		util.DebugError("UserGetByIDs:", err)
		return nil
	}

	return users
}

// 获取所有学生数据(谨慎使用)
func UserGetAll(session *xorm.Session) []Student {
	var list []Student
	// Find 直接把查询结果填充到切片
	err := session.Find(&list)
	if err != nil {
		// 出错返回空
		util.DebugError("UserGetAll:", err)
		return nil
	}
	return list
}
