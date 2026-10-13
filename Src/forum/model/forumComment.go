package model

import "time"

type ForumComment struct {
	Indices                    int64     `json:"indices" xorm:"'indices' pk autoincr"`
	IndicesForForumAnnoucement int64     `json:"indicesForForumAnnoucement" xorm:"notnull 'indicesForForumAnnoucement'"`
	ID                         string    `json:"id" xorm:"varchar(50) notnull 'id'"`
	Content                    string    `json:"content" xorm:"text 'content'"`
	SubmitTime                 time.Time `json:"submittime"`
}

func (p ForumComment) TableName() string {
	return "ForumComment"
}
