package dao

import (
	database "newaoe/Src/databse"
	"newaoe/Src/user/model"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func UserResetAvatar(session *xorm.Session, id string, avatarFile string) (bool, error) {
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
func UserExist(session *xorm.Session, id string) (bool, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	has, err := session.Where("id = ?", id).Exist(&model.Student{})
	if (!has) || (err != nil) {
		return false, err
	}
	return true, nil
}

// 判断邮箱是否已经使用
func EmailExist(session *xorm.Session, email string) (bool, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	has, err := session.Where("email = ?", email).Exist(&model.Student{})
	if (!has) || (err != nil) {
		return false, err
	}
	return true, nil
}

// 添加用户
func UserAdd(session *xorm.Session, userInfo model.Student) (bool, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data := &model.Student{
		Id:         userInfo.Id,
		Password:   userInfo.Password,
		Email:      userInfo.Email,
		Vip:        userInfo.Vip,
		RegistDate: util.UTC_Time(),
		Avatar:     util.GetRandomAvatar(),
	}
	_, err := session.Insert(data)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 添加用户
func UserAddByStudentInfo(session *xorm.Session, info model.Student) (bool, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	_, err := session.Insert(info)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 更新用户数据
func UserUpdate(session *xorm.Session, id string, user model.Student) (bool, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	_, err := session.Where("id = ?", id).Update(&user)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 获取 id like 的所有用户(不包含end)
func UserGetLikeId(session *xorm.Session, id_like string, beg int, end int) ([]model.Student, error) {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	var users []model.Student

	err := session.Where("id LIKE ?", "%"+id_like+"%").Limit(end-beg, beg).Find(&users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// 获取指定用户的数据
func UserGet(session *xorm.Session, id string) (*model.Student, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var user model.Student
	has, err := session.Where("id = ? ", id).Get(&user)
	if err != nil || !has {
		return nil, err
	}
	return &user, nil
}

// 获取指定用户的数据
func UserGetByEmail(session *xorm.Session, email string) (*model.Student, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var user model.Student
	has, err := session.Where("email = ? ", email).Get(&user)
	if err != nil || !has {
		return nil, err
	}
	return &user, nil
}

// 获取指定用户的数据
func UserGetByIdOrEmail(session *xorm.Session, id_or_email string) (*model.Student, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	var user model.Student
	has, err := session.Where("id = ? OR email = ?", id_or_email, id_or_email).Get(&user)
	if err != nil || !has {
		return nil, err
	}
	return &user, nil
}

// 根据注册时间升序排序（不包括end)
func UserGetByRangeOrderByRegistData(session *xorm.Session, beg int, end int) ([]model.Student, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	if beg < 0 || beg >= end {
		return nil, util.NewError("beg<0 or beg>=end!")
	}
	var list []model.Student
	// 按regist_date升序
	// LIMIT beg, end-beg  等价于 offset beg limit count
	err := session.Asc("registDate").Limit(end-beg, beg).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 获取指定多个用户的数据
func UserGetByIDs(session *xorm.Session, ids []string) ([]model.Student, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//

	var users []model.Student

	if len(ids) == 0 {
		return users, nil
	}

	err := session.In("id", ids).Find(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// 获取所有学生数据(谨慎使用)
func UserGetAll(session *xorm.Session) ([]model.Student, error) {
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
		return nil, err
	}
	return list, nil
}

// 获取用户密码(编码后的密码)
func UserGetPassword(session *xorm.Session, id string) (string, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data, err := UserGet(session, id)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	return data.Password, nil
}

func UserGetEmail(session *xorm.Session, id string) (string, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data, err := UserGet(session, id)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	return data.Email, nil
}

func UserGetAvatar(session *xorm.Session, id string) (string, error) {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data, err := UserGet(session, id)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	return data.Avatar, nil
}

func UserResetPassword(session *xorm.Session, id string, new_password string) (bool, error) {
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

// 获取 [beg, end)，不包含 end
func UserGetRangeIDs(session *xorm.Session, beg int, end int) ([]string, error) {
	if beg < 0 || end <= beg {
		return nil, util.NewError("beg<0 or beg>=end!")
	}

	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}

	var ids []string
	err := session.
		Table(model.Student{}.TableName()).
		Cols("id").
		OrderBy("id ASC").
		Limit(end-beg, beg).
		Find(&ids)

	if err != nil {
		return nil, err
	}

	return ids, nil
}
