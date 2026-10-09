package dao

import (
	"newaoe/Src/codeAssessmentSubmit/model"
	database "newaoe/Src/databse"
	"newaoe/Src/util"

	"xorm.io/xorm"
)

func TeacherAdd(session *xorm.Session, teacher string) error {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	data := model.Teacher{
		Name: teacher,
	}
	_, err := session.Insert(&data)
	if err != nil {
		return err
	}
	return nil
}

func TeacherGetAll(session *xorm.Session) []model.Teacher {
	//
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//从数据库查询数据
	var teachers []model.Teacher
	err := session.Find(&teachers)
	if err != nil {
		util.DebugError("TeacherGetAll:", err)
		return nil
	}
	//返回数据
	return teachers
}

func TeacherExist(session *xorm.Session, name string) bool {
	if session == nil {
		session = database.NewSession()
		defer session.Close()
	}
	//
	teacher := TeacherGetAll(session)
	if len(teacher) == 0 {
		return false
	}
	//
	for _, id := range teacher {
		if id.Name == name {
			return true
		}
	}
	//
	return false
}
