package dao

import "newaoe/util"

//本次授课老师
type Teacher struct {
	Name string `json:"name" xorm:"name"`
}

func (t Teacher) TableName() string {
	return "Teacher"
}

func TeacherGetAll() ([]Teacher, error) {
	//
	var ret []Teacher
	err := DB.Find(&ret)
	if err != nil {
		util.DebugError("TeacherGetAll:", err)
		return nil, err
	}
	return ret, nil
}
