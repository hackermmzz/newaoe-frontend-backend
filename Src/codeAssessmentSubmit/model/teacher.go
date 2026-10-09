package model

// 本次授课老师
type Teacher struct {
	Name string `json:"name" xorm:"name"`
}

func (t Teacher) TableName() string {
	return "Teacher"
}
