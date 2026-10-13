package model

import "time"

type ForumAnnouncement struct {
	Indices     int       `json:"indices" xorm:"'indices' pk autoincr"`
	ID          string    `json:"id" xorm:"'id'" `
	Title       string    `json:"title" xorm:"'title'" `
	ContentPath string    `json:"contentPath" xorm:"'contentPath'" `
	Summary     string    `json:"summary" xorm:"'summary'" `
	SubmitTime  time.Time `json:"submitTime" xorm:"'submitTime'" `
}

func (p ForumAnnouncement) TableName() string {
	return "ForumAnnouncement"
}
